// Copyright 2020 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package registration

import (
	"context"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/common/logging"
	"go.chromium.org/luci/gae/service/datastore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ufspb "infra/unifiedfleet/api/v1/models"
	ufsds "infra/unifiedfleet/app/model/datastore"
	"infra/unifiedfleet/app/util"
)

// SwitchKind is the datastore entity kind Switch.
const SwitchKind string = "Switch"

// SwitchEntity is a datastore entity that tracks switch.
type SwitchEntity struct {
	_kind string   `gae:"$kind,Switch"`
	ID    string   `gae:"$id"`
	Lab   string   `gae:"lab"` // deprecated
	Zone  string   `gae:"zone"`
	Rack  string   `gae:"rack"`
	Tags  []string `gae:"tags"`
	State string   `gae:"state"`
	// ufspb.Switch cannot be directly used as it contains pointer (timestamp).
	Switch []byte `gae:",noindex"`
}

// GetProto returns the unmarshaled switch.
func (e *SwitchEntity) GetProto() (proto.Message, error) {
	var p ufspb.Switch
	if err := proto.Unmarshal(e.Switch, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func newSwitchEntity(ctx context.Context, pm proto.Message) (ufsds.FleetEntity, error) {
	p := pm.(*ufspb.Switch)
	if p.GetName() == "" {
		return nil, errors.Reason("Empty Switch ID").Err()
	}
	s, err := proto.Marshal(p)
	if err != nil {
		return nil, errors.Annotate(err, "fail to marshal Switch %s", p).Err()
	}
	return &SwitchEntity{
		ID:     p.GetName(),
		Switch: s,
		Zone:   p.GetZone(),
		Rack:   p.GetRack(),
		Tags:   p.GetTags(),
		State:  p.GetResourceState().String(),
	}, nil
}

// CreateSwitch creates a new switch in datastore.
func CreateSwitch(ctx context.Context, s *ufspb.Switch) (*ufspb.Switch, error) {
	return putSwitch(ctx, s, false)
}

// UpdateSwitch updates switch in datastore.
func UpdateSwitch(ctx context.Context, s *ufspb.Switch) (*ufspb.Switch, error) {
	return putSwitch(ctx, s, true)
}

// GetSwitch returns switch for the given id from datastore.
func GetSwitch(ctx context.Context, id string) (*ufspb.Switch, error) {
	pm, err := ufsds.Get(ctx, &ufspb.Switch{Name: id}, newSwitchEntity)
	if err == nil {
		return pm.(*ufspb.Switch), err
	}
	return nil, err
}

func getSwitchID(pm proto.Message) string {
	p := pm.(*ufspb.Switch)
	return p.GetName()
}

// BatchGetSwitches returns a batch of switches from datastore.
func BatchGetSwitches(ctx context.Context, ids []string) ([]*ufspb.Switch, error) {
	protos := make([]proto.Message, len(ids))
	for i, n := range ids {
		protos[i] = &ufspb.Switch{Name: n}
	}
	pms, err := ufsds.BatchGet(ctx, protos, newSwitchEntity, getSwitchID)
	if err != nil {
		return nil, err
	}
	res := make([]*ufspb.Switch, len(pms))
	for i, pm := range pms {
		res[i] = pm.(*ufspb.Switch)
	}
	return res, nil
}

// QuerySwitchByPropertyName query's Switch Entity in the datastore
//
// If keysOnly is true, then only key field is populated in returned switches
func QuerySwitchByPropertyName(ctx context.Context, propertyName, id string, keysOnly bool) ([]*ufspb.Switch, error) {
	q := datastore.NewQuery(SwitchKind).KeysOnly(keysOnly).FirestoreMode(true)
	var entities []*SwitchEntity
	if err := datastore.GetAll(ctx, q.Eq(propertyName, id), &entities); err != nil {
		logging.Errorf(ctx, "Failed to query from datastore: %s", err)
		return nil, status.Errorf(codes.Internal, ufsds.InternalError)
	}
	if len(entities) == 0 {
		logging.Debugf(ctx, "No switches found for the query: %s", id)
		return nil, nil
	}
	switches := make([]*ufspb.Switch, 0, len(entities))
	for _, entity := range entities {
		if keysOnly {
			s := &ufspb.Switch{
				Name: entity.ID,
			}
			switches = append(switches, s)
		} else {
			pm, perr := entity.GetProto()
			if perr != nil {
				logging.Errorf(ctx, "Failed to unmarshal proto: %s", perr)
				continue
			}
			switches = append(switches, pm.(*ufspb.Switch))
		}
	}
	return switches, nil
}

// ListSwitches lists the switches
//
// Does a query over switch entities. Returns up to pageSize entities, plus non-nil cursor (if
// there are more results). pageSize must be positive.
func ListSwitches(ctx context.Context, pageSize int32, pageToken string, filterMap map[string][]interface{}, keysOnly bool) (res []*ufspb.Switch, nextPageToken string, err error) {
	q, err := ufsds.ListQuery(ctx, SwitchKind, pageSize, pageToken, filterMap, keysOnly)
	if err != nil {
		return nil, "", err
	}
	var nextCur datastore.Cursor
	err = datastore.Run(ctx, q, func(ent *SwitchEntity, cb datastore.CursorCB) error {
		if keysOnly {
			s := &ufspb.Switch{
				Name: ent.ID,
			}
			res = append(res, s)
		} else {
			pm, err := ent.GetProto()
			if err != nil {
				logging.Errorf(ctx, "Failed to UnMarshal: %s", err)
				return nil
			}
			res = append(res, pm.(*ufspb.Switch))
		}
		if len(res) >= int(pageSize) {
			if nextCur, err = cb(); err != nil {
				return err
			}
			return datastore.Stop
		}
		return nil
	})
	if err != nil {
		logging.Errorf(ctx, "Failed to List Switches %s", err)
		return nil, "", status.Errorf(codes.Internal, ufsds.InternalError)
	}
	if nextCur != nil {
		nextPageToken = nextCur.String()
	}
	return
}

// DeleteSwitch deletes the switch in datastore
func DeleteSwitch(ctx context.Context, id string) error {
	return ufsds.Delete(ctx, &ufspb.Switch{Name: id}, newSwitchEntity)
}

func putSwitch(ctx context.Context, s *ufspb.Switch, update bool) (*ufspb.Switch, error) {
	s.UpdateTime = ptypes.TimestampNow()
	pm, err := ufsds.Put(ctx, s, newSwitchEntity, update)
	if err == nil {
		return pm.(*ufspb.Switch), err
	}
	return nil, err
}

// ImportSwitches creates or updates a batch of switches in datastore
func ImportSwitches(ctx context.Context, switches []*ufspb.Switch) (*ufsds.OpResults, error) {
	protos := make([]proto.Message, len(switches))
	utime := ptypes.TimestampNow()
	for i, m := range switches {
		m.UpdateTime = utime
		protos[i] = m
	}
	return ufsds.Insert(ctx, protos, newSwitchEntity, true, true)
}

// BatchUpdateSwitches updates switches in datastore.
//
// This is a non-atomic operation and doesnt check if the object already exists before
// update. Must be used within a Transaction where objects are checked before update.
// Will lead to partial updates if not used in a transaction.
func BatchUpdateSwitches(ctx context.Context, switches []*ufspb.Switch) ([]*ufspb.Switch, error) {
	return putAllSwitch(ctx, switches, true)
}

func putAllSwitch(ctx context.Context, switches []*ufspb.Switch, update bool) ([]*ufspb.Switch, error) {
	protos := make([]proto.Message, len(switches))
	updateTime := ptypes.TimestampNow()
	for i, s := range switches {
		s.UpdateTime = updateTime
		protos[i] = s
	}
	_, err := ufsds.PutAll(ctx, protos, newSwitchEntity, update)
	if err == nil {
		return switches, err
	}
	return nil, err
}

func queryAllSwitch(ctx context.Context) ([]ufsds.FleetEntity, error) {
	var entities []*SwitchEntity
	q := datastore.NewQuery(SwitchKind)
	if err := datastore.GetAll(ctx, q, &entities); err != nil {
		return nil, err
	}
	fe := make([]ufsds.FleetEntity, len(entities))
	for i, e := range entities {
		fe[i] = e
	}
	return fe, nil
}

// GetAllSwitches returns all switches in datastore.
func GetAllSwitches(ctx context.Context) (*ufsds.OpResults, error) {
	return ufsds.GetAll(ctx, queryAllSwitch)
}

// DeleteSwitches deletes a batch of switches
func DeleteSwitches(ctx context.Context, resourceNames []string) *ufsds.OpResults {
	protos := make([]proto.Message, len(resourceNames))
	for i, m := range resourceNames {
		protos[i] = &ufspb.Switch{
			Name: m,
		}
	}
	return ufsds.DeleteAll(ctx, protos, newSwitchEntity)
}

// BatchDeleteSwitches deletes switches in datastore.
//
// This is a non-atomic operation. Must be used within a transaction.
// Will lead to partial deletes if not used in a transaction.
func BatchDeleteSwitches(ctx context.Context, ids []string) error {
	protos := make([]proto.Message, len(ids))
	for i, id := range ids {
		protos[i] = &ufspb.Switch{Name: id}
	}
	return ufsds.BatchDelete(ctx, protos, newSwitchEntity)
}

// GetSwitchIndexedFieldName returns the index name
func GetSwitchIndexedFieldName(input string) (string, error) {
	var field string
	input = strings.TrimSpace(input)
	switch strings.ToLower(input) {
	case util.ZoneFilterName:
		field = "zone"
	case util.RackFilterName:
		field = "rack"
	case util.TagFilterName:
		field = "tags"
	case util.StateFilterName:
		field = "state"
	default:
		return "", status.Errorf(codes.InvalidArgument, "Invalid field name %s - field name for Switch are zone/rack/tag/state", input)
	}
	return field, nil
}
