// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package internal

import (
	"bytes"
	"crypto/sha1"
	"errors"

	"github.com/golang/protobuf/proto"

	"infra/tools/cipd/internal/messages"
)

// MarshalWithSHA1 serializes proto message to bytes, calculates SHA1 checksum
// of it, and returns serialized envelope that contains both. UnmarshalWithSHA1
// can then be used to verify SHA1 and deserialized the original object.
func MarshalWithSHA1(pm proto.Message) ([]byte, error) {
	blob, err := proto.Marshal(pm)
	if err != nil {
		return nil, err
	}
	sum := sha1.Sum(blob)
	envelope := messages.BlobWithSHA1{Blob: blob, Sha1: sum[:]}
	return proto.Marshal(&envelope)
}

// UnmarshalWithSHA1Check is reverse of MarshalWithSHA1. It checks SHA1 checksum
// and deserializes the object if it matches the blob.
func UnmarshalWithSHA1(buf []byte, pm proto.Message) error {
	envelope := messages.BlobWithSHA1{}
	if err := proto.Unmarshal(buf, &envelope); err != nil {
		return err
	}
	sum := sha1.Sum(envelope.GetBlob())
	if !bytes.Equal(sum[:], envelope.GetSha1()) {
		return errors.New("check sum of tag cache file is invalid")
	}
	return proto.Unmarshal(envelope.GetBlob(), pm)
}
