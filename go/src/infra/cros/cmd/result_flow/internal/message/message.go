// Copyright 2020 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package message

import (
	"context"
	"encoding/json"
	"fmt"
	"infra/libs/skylab/request"
	"strconv"

	"go.chromium.org/chromiumos/infra/proto/go/test_platform/result_flow"
	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/common/logging"

	pubsub "cloud.google.com/go/pubsub/apiv1"
	"google.golang.org/api/option"
	pubsubpb "google.golang.org/genproto/googleapis/pubsub/v1"
)

const (
	// BuildIDKeyName is the key name to store Build ID in message attributes.
	BuildIDKeyName = "build_id"
	// DefaultMaxReceivingMessage is default max message a single flow could handle.
	DefaultMaxReceivingMessage = 15
)

// PublishBuildID publishes a Build ID to PubSub.
func PublishBuildID(ctx context.Context, bID int64, conf *result_flow.PubSubConfig, opts ...option.ClientOption) error {
	c, err := pubsub.NewPublisherClient(ctx, opts...)
	if err != nil {
		return fmt.Errorf("failed to create publisher client: %v", err)
	}
	defer c.Close()

	_, err = c.Publish(ctx, &pubsubpb.PublishRequest{
		Topic: fmt.Sprintf("projects/%s/topics/%s", conf.Project, conf.Topic),
		Messages: []*pubsubpb.PubsubMessage{
			genPublishRequest(bID),
		},
	})
	return err
}

func genPublishRequest(bID int64) *pubsubpb.PubsubMessage {
	return &pubsubpb.PubsubMessage{
		Attributes: map[string]string{
			BuildIDKeyName: fmt.Sprintf("%d", bID),
		},
	}
}

// Client defines an interface used to interact with pubsub
type Client interface {
	PullMessages(context.Context) ([]*pubsubpb.ReceivedMessage, error)
	AckMessages(context.Context, []*pubsubpb.ReceivedMessage) error
	Close() error
}

type messageClient struct {
	client       *pubsub.SubscriberClient
	subscription string
	maxMessages  int32
}

// NewClient creates a messageClient for PubSub subscriber.
func NewClient(c context.Context, conf *result_flow.PubSubConfig, opts ...option.ClientOption) (Client, error) {
	client, err := pubsub.NewSubscriberClient(c, opts...)
	if err != nil {
		return nil, err
	}
	maxMessages := conf.GetMaxReceivingMessages()
	if maxMessages == 0 {
		maxMessages = int32(DefaultMaxReceivingMessage)
	}
	return &messageClient{
		client:       client,
		subscription: fmt.Sprintf("projects/%s/subscriptions/%s", conf.Project, conf.Subscription),
		maxMessages:  maxMessages,
	}, nil
}

// PullMessages fetches messages from Pub/Sub.
func (m *messageClient) PullMessages(c context.Context) ([]*pubsubpb.ReceivedMessage, error) {
	req := pubsubpb.PullRequest{
		Subscription: m.subscription,
		MaxMessages:  m.maxMessages,
	}
	res, err := m.client.Pull(c, &req)
	if err != nil {
		return nil, errors.Annotate(err, "failed to pull messages").Err()
	}
	return res.ReceivedMessages, nil
}

// AckMessages acknowledges a list of messages.
func (m *messageClient) AckMessages(c context.Context, msgs []*pubsubpb.ReceivedMessage) error {
	if len(msgs) == 0 {
		return nil
	}
	AckIDs := make([]string, len(msgs))
	for i, msg := range msgs {
		AckIDs[i] = msg.AckId
	}
	return m.client.Acknowledge(c, &pubsubpb.AcknowledgeRequest{
		Subscription: m.subscription,
		AckIds:       AckIDs,
	})
}

// Close terminates the subscriber client.
func (m *messageClient) Close() error {
	return m.client.Close()
}

// ExtractBuildIDMap generates a map with key of Build ID and the value of parent UID.
func ExtractBuildIDMap(ctx context.Context, msgs []*pubsubpb.ReceivedMessage) map[int64]*pubsubpb.ReceivedMessage {
	m := make(map[int64]*pubsubpb.ReceivedMessage, len(msgs))
	for _, msg := range msgs {
		bID, err := extractBuildID(msg)
		if err != nil {
			logging.Errorf(ctx, "Failed to extract build ID, err: %v", err)
			continue
		}
		m[bID] = msg
	}
	return m
}

// ExtractParentUID extracts the parent request UID from the pubsub message.
func ExtractParentUID(msg *pubsubpb.ReceivedMessage) (string, error) {
	msgBody := struct {
		UserData string `json:"user_data"`
	}{}
	if err := json.Unmarshal(msg.GetMessage().GetData(), &msgBody); err != nil {
		return "", errors.Annotate(err, "could not parse pubsub message data").Err()
	}
	msgPayload := request.MessagePayload{}
	if err := json.Unmarshal([]byte(msgBody.UserData), &msgPayload); err != nil {
		return "", errors.Annotate(err, "could not extract Parent UID").Err()
	}
	return msgPayload.ParentRequestUID, nil
}

func extractBuildID(msg *pubsubpb.ReceivedMessage) (int64, error) {
	bID, err := strconv.ParseInt(msg.Message.Attributes[BuildIDKeyName], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Failed to parse build ID from: %s", msg.Message.Attributes[BuildIDKeyName])
	}
	if bID == 0 {
		return 0, fmt.Errorf("Build ID can not be 0")
	}
	return bID, nil
}
