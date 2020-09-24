// Copyright 2020 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package controller

import (
	"context"
	"strings"

	"github.com/golang/protobuf/proto"
	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/common/logging"
	"go.chromium.org/luci/gae/service/datastore"
	crimsonconfig "go.chromium.org/luci/machine-db/api/config/v1"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"infra/libs/cros/git"
	"infra/libs/cros/sheet"
	ufspb "infra/unifiedfleet/api/v1/proto"
	ufsAPI "infra/unifiedfleet/api/v1/rpc"
	"infra/unifiedfleet/app/config"
	"infra/unifiedfleet/app/model/configuration"
	ufsds "infra/unifiedfleet/app/model/datastore"
	"infra/unifiedfleet/app/model/inventory"
	"infra/unifiedfleet/app/util"
)

// CreateVlan creates a new vlan in datastore.
func CreateVlan(ctx context.Context, vlan *ufspb.Vlan) (*ufspb.Vlan, error) {
	var ips []*ufspb.IP
	var length int
	var err error
	var freeStartIP, freeEndIP string
	f := func(ctx context.Context) error {
		hc := getVlanHistoryClient(vlan)
		if err := validateCreateVlan(ctx, vlan); err != nil {
			return errors.Annotate(err, "CreateVlan - validation failed").Err()
		}

		ips, length, freeStartIP, freeEndIP, err = util.ParseVlan(vlan.GetName(), vlan.GetVlanAddress())
		if err != nil {
			return errors.Annotate(err, "CreateVlan").Err()
		}
		vlan.CapacityIp = int32(length)
		vlan.ResourceState = ufspb.State_STATE_SERVING
		vlan.VlanNumber = util.GetSuffixAfterSeparator(vlan.GetName(), ":")
		vlan.FreeStartIpv4Str = freeStartIP
		vlan.FreeEndIpv4Str = freeEndIP

		if _, err = configuration.BatchUpdateVlans(ctx, []*ufspb.Vlan{vlan}); err != nil {
			return err
		}
		hc.stUdt.updateStateHelper(ctx, ufspb.State_STATE_SERVING)
		hc.LogVLANChanges(nil, vlan)
		return hc.SaveChangeEvents(ctx)
	}

	if err := datastore.RunInTransaction(ctx, f, nil); err != nil {
		return nil, errors.Annotate(err, "CreateVlan - unable to create vlan %s", vlan.Name).Err()
	}
	// Cannot update more than 500 entities in one transaction
	logging.Debugf(ctx, "Updating %d ips", len(ips))
	for i := 0; ; i += util.OperationPageSize {
		end := util.Min(i+util.OperationPageSize, len(ips))
		configuration.BatchUpdateIPs(ctx, ips[i:end])
		if i+util.OperationPageSize >= len(ips) {
			break
		}
	}
	return vlan, nil
}

// UpdateVlan updates vlan in datastore.
func UpdateVlan(ctx context.Context, vlan *ufspb.Vlan, mask *field_mask.FieldMask) (*ufspb.Vlan, error) {
	f := func(ctx context.Context) error {
		hc := getVlanHistoryClient(vlan)
		if err := validateUpdateVlan(ctx, vlan, mask); err != nil {
			return errors.Annotate(err, "UpdateVlan - validation failed").Err()
		}

		oldVlan, err := configuration.GetVlan(ctx, vlan.GetName())
		if err != nil {
			return errors.Annotate(err, "UpdateVlan - fail to get old vlan").Err()
		}
		oldVlanCopy := proto.Clone(oldVlan).(*ufspb.Vlan)

		// Copy the not-allowed change fields
		vlan.VlanAddress = oldVlan.GetVlanAddress()
		vlan.CapacityIp = oldVlan.GetCapacityIp()
		vlan.VlanNumber = oldVlan.GetVlanNumber()

		if err := validateReservedIPs(ctx, vlan); err != nil {
			return err
		}

		// Partial update by field mask
		if mask != nil && len(mask.Paths) > 0 {
			vlan, err = processVlanUpdateMask(oldVlan, vlan, mask)
			if err != nil {
				return errors.Annotate(err, "UpdateVlan - processing update mask failed").Err()
			}
		}

		// update state
		if err := hc.stUdt.updateStateHelper(ctx, vlan.GetResourceState()); err != nil {
			return errors.Annotate(err, "Fail to update state to vlan %s", vlan.GetName()).Err()
		}

		if err := updateIPTable(ctx, vlan.GetName(), oldVlanCopy.GetReservedIps(), vlan.GetReservedIps()); err != nil {
			return err
		}

		if _, err := configuration.BatchUpdateVlans(ctx, []*ufspb.Vlan{vlan}); err != nil {
			return errors.Annotate(err, "UpdateVlan - unable to batch update vlan %s", vlan.Name).Err()
		}
		hc.LogVLANChanges(oldVlanCopy, vlan)
		return hc.SaveChangeEvents(ctx)
	}

	if err := datastore.RunInTransaction(ctx, f, nil); err != nil {
		return nil, errors.Annotate(err, "UpdateVlan - unable to update vlan %s", vlan.Name).Err()
	}
	return vlan, nil
}

// GetVlan returns vlan for the given id from datastore.
func GetVlan(ctx context.Context, id string) (*ufspb.Vlan, error) {
	return configuration.GetVlan(ctx, id)
}

// BatchGetVlans returns a batch of vlans from datastore.
func BatchGetVlans(ctx context.Context, ids []string) ([]*ufspb.Vlan, error) {
	return configuration.BatchGetVlans(ctx, ids)
}

// ListVlans lists the vlans
func ListVlans(ctx context.Context, pageSize int32, pageToken, filter string, keysOnly bool) ([]*ufspb.Vlan, string, error) {
	var filterMap map[string][]interface{}
	var err error
	if filter != "" {
		filterMap, err = getFilterMap(filter, configuration.GetVlanIndexedFieldName)
		if err != nil {
			return nil, "", errors.Annotate(err, "Failed to read filter for listing vlans").Err()
		}
	}
	filterMap = resetStateFilter(filterMap)
	return configuration.ListVlans(ctx, pageSize, pageToken, filterMap, keysOnly)
}

// DeleteVlan deletes the vlan in datastore
//
// For referential data intergrity,
// Delete if this Vlan is not referenced by other resources in the datastore.
// If there are any references, delete will be rejected and an error will be returned.
func DeleteVlan(ctx context.Context, id string) error {
	f := func(ctx context.Context) error {
		hc := getVlanHistoryClient(&ufspb.Vlan{Name: id})
		hc.LogVLANChanges(nil, &ufspb.Vlan{Name: id})
		err := validateDeleteVlan(ctx, id)
		if err != nil {
			return errors.Annotate(err, "DeleteVlan - validation failed").Err()
		}

		if err := configuration.DeleteVlan(ctx, id); err != nil {
			return err
		}
		return hc.SaveChangeEvents(ctx)
	}

	if err := datastore.RunInTransaction(ctx, f, nil); err != nil {
		return errors.Annotate(err, "DeleteVlan - unable to delete vlan %s", id).Err()
	}
	ips, _ := configuration.QueryIPByPropertyName(ctx, map[string]string{
		"vlan": id,
	})
	// Cannot update more than 500 entities in one transaction
	if len(ips) > 0 {
		logging.Debugf(ctx, "deleting %d ips for vlan %s", len(ips), id)
		ipIDs := ufsAPI.ParseResources(ips, "Id")
		for i := 0; ; i += util.OperationPageSize {
			end := util.Min(i+util.OperationPageSize, len(ipIDs))
			configuration.BatchDeleteIPs(ctx, ipIDs[i:end])
			if i+util.OperationPageSize >= len(ipIDs) {
				break
			}
		}
	}
	return nil
}

// ImportVlans implements the logic of importing vlans and related info to backend storage.
//
// The function will return:
//      * all of the results of the operations that already run
//      * the first error that it meets
//
// The function will stop at the very first error.
func ImportVlans(ctx context.Context, vlans []*crimsonconfig.VLAN, pageSize int) (*ufsds.OpResults, error) {
	IPs := make([]*ufspb.IP, 0)
	vs := make([]*ufspb.Vlan, len(vlans))
	for i, vlan := range vlans {
		vlanName := util.GetBrowserLabName(util.Int64ToStr(vlan.GetId()))
		ips, length, freeStartIP, freeEndIP, err := util.ParseVlan(vlanName, vlan.GetCidrBlock())
		if err != nil {
			return nil, err
		}
		for _, ip := range ips {
			IPs = append(IPs, ip)
		}
		vs[i] = &ufspb.Vlan{
			Name:             vlanName,
			Description:      vlan.GetAlias(),
			CapacityIp:       int32(length),
			VlanAddress:      vlan.GetCidrBlock(),
			ResourceState:    util.ToState(vlan.GetState()),
			FreeStartIpv4Str: freeStartIP,
			FreeEndIpv4Str:   freeEndIP,
			VlanNumber:       util.GetSuffixAfterSeparator(vlanName, ":"),
		}
	}
	deleteNonExistingVlans(ctx, vs, pageSize)
	allRes := make(ufsds.OpResults, 0)
	logging.Infof(ctx, "Importing %d vlans", len(vs))
	for i := 0; ; i += pageSize {
		end := util.Min(i+pageSize, len(vs))
		res, err := configuration.ImportVlans(ctx, vs[i:end])
		allRes = append(allRes, *res...)
		if err != nil {
			return &allRes, err
		}
		if i+pageSize >= len(vs) {
			break
		}
	}

	deleteInvalidIPs(ctx, pageSize)
	logging.Infof(ctx, "Importing %d ips", len(IPs))
	for i := 0; ; i += pageSize {
		end := util.Min(i+pageSize, len(IPs))
		res, err := configuration.ImportIPs(ctx, IPs[i:end])
		allRes = append(allRes, *res...)
		if err != nil {
			return &allRes, err
		}
		if i+pageSize >= len(IPs) {
			break
		}
	}
	return &allRes, nil
}

func deleteInvalidIPs(ctx context.Context, pageSize int) {
	resp, err := configuration.GetAllIPs(ctx)
	if err != nil {
		logging.Debugf(ctx, "Fail to get all ips: %s", err.Error())
		return
	}
	var toDeleteIP []string
	for _, sr := range resp.Passed() {
		s := sr.Data.(*ufspb.IP)
		if s.GetIpv4Str() == "" {
			toDeleteIP = append(toDeleteIP, s.GetId())
		}
	}
	logging.Infof(ctx, "Deleting %d invalid ips ", len(toDeleteIP))
	deleteByPage(ctx, toDeleteIP, pageSize, configuration.DeleteIPs)
}

func deleteNonExistingVlans(ctx context.Context, vlans []*ufspb.Vlan, pageSize int) (*ufsds.OpResults, error) {
	resMap := make(map[string]bool)
	for _, r := range vlans {
		resMap[r.GetName()] = true
	}
	resp, err := configuration.GetAllVlans(ctx)
	if err != nil {
		return nil, err
	}
	var toDelete []string
	var toDeleteIP []string
	for _, sr := range resp.Passed() {
		s := sr.Data.(*ufspb.Vlan)
		if util.IsInBrowserZone(s.GetName()) || strings.Contains(s.GetName(), "browser-lab") {
			if _, ok := resMap[s.GetName()]; !ok {
				toDelete = append(toDelete, s.GetName())
				ips, err := configuration.QueryIPByPropertyName(ctx, map[string]string{"vlan": s.GetName()})
				if err != nil {
					return nil, err
				}
				for _, ip := range ips {
					toDeleteIP = append(toDeleteIP, ip.GetId())
				}
			}
		}
	}

	logging.Infof(ctx, "Deleting %d non-existing ips ", len(toDeleteIP))
	deleteByPage(ctx, toDeleteIP, pageSize, configuration.DeleteIPs)
	logging.Infof(ctx, "Deleting %d non-existing vlans ", len(toDelete))
	return deleteByPage(ctx, toDelete, pageSize, configuration.DeleteVlans), nil
}

// ImportOSVlans imports the logic of parse and save network infos.
func ImportOSVlans(ctx context.Context, sheetClient sheet.ClientInterface, gitClient git.ClientInterface, pageSize int) (*ufsds.OpResults, error) {
	networkCfg := config.Get(ctx).GetCrosNetworkConfig()
	allVlans := make([]*ufspb.Vlan, 0)
	allIPs := make([]*ufspb.IP, 0)
	allDhcps := make([]*ufspb.DHCPConfig, 0)

	for _, cfg := range networkCfg.GetCrosNetworkTopology() {
		logging.Debugf(ctx, "########### Parse %s ###########", cfg.GetName())
		resp, err := sheetClient.Get(ctx, cfg.GetSheetId(), []string{"VLANs and Netblocks"})
		if err != nil {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		topology, dupcatedVlan := util.ParseATLTopology(resp)
		logging.Debugf(ctx, "Topology length %d", len(topology))
		logging.Debugf(ctx, "Duplicated vlans found in topology:")
		logVlans(ctx, dupcatedVlan)
		conf, err := gitClient.GetFile(ctx, cfg.GetRemotePath())
		if err != nil {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		parsed, err := util.ParseOSDhcpdConf(conf, topology)
		if err != nil {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		logging.Debugf(ctx, "Duplicated vlans found in dhcp conf file:")
		logVlans(ctx, parsed.DuplicatedVlans)
		logging.Debugf(ctx, "Vlans not existing in pre-defined topology:")
		logVlans(ctx, parsed.MismatchedVlans)
		logging.Debugf(ctx, "Invalid dhcps without vlan:")
		logDHCPs(ctx, parsed.DHCPsWithoutVlan)
		logging.Debugf(ctx, "Duplicated ips found in dhcp conf file:")
		logIPs(ctx, parsed.DuplicatedIPs)

		logging.Debugf(ctx, "Get %d vlans, %d ips, %d dhcps for %s", len(parsed.ValidVlans), len(parsed.ValidIPs), len(parsed.ValidDHCPs), cfg.GetName())
		allVlans = append(allVlans, parsed.ValidVlans...)
		allIPs = append(allIPs, parsed.ValidIPs...)
		allDhcps = append(allDhcps, parsed.ValidDHCPs...)
	}

	allRes := make(ufsds.OpResults, 0)
	logging.Debugf(ctx, "Importing %d vlans", len(allVlans))
	for i := 0; ; i += pageSize {
		end := util.Min(i+pageSize, len(allVlans))
		res, err := configuration.ImportVlans(ctx, allVlans[i:end])
		allRes = append(allRes, *res...)
		if err != nil {
			return &allRes, err
		}
		if i+pageSize >= len(allVlans) {
			break
		}
	}

	logging.Debugf(ctx, "Importing %d ips", len(allIPs))
	for i := 0; ; i += pageSize {
		end := util.Min(i+pageSize, len(allIPs))
		res, err := configuration.ImportIPs(ctx, allIPs[i:end])
		allRes = append(allRes, *res...)
		if err != nil {
			return &allRes, err
		}
		if i+pageSize >= len(allIPs) {
			break
		}
	}

	logging.Debugf(ctx, "Importing %d ips", len(allDhcps))
	for i := 0; ; i += pageSize {
		end := util.Min(i+pageSize, len(allDhcps))
		res, err := configuration.ImportDHCPConfigs(ctx, allDhcps[i:end])
		allRes = append(allRes, *res...)
		if err != nil {
			return &allRes, err
		}
		if i+pageSize >= len(allDhcps) {
			break
		}
	}
	return &allRes, nil
}

func logVlans(ctx context.Context, vlans []*ufspb.Vlan) {
	if vlans != nil && len(vlans) > 0 {
		for _, v := range vlans {
			logging.Debugf(ctx, "\tVlan %s (%s)", v.GetName(), v.GetVlanAddress())
		}
		return
	}
	logging.Debugf(ctx, "\tNot found")
}

func logDHCPs(ctx context.Context, dhcps []*ufspb.DHCPConfig) {
	if dhcps != nil && len(dhcps) > 0 {
		for _, v := range dhcps {
			logging.Debugf(ctx, "\tHost %s (%s)", v.GetHostname(), v.GetIp())
		}
		return
	}
	logging.Debugf(ctx, "\tNot found")

}

func logIPs(ctx context.Context, ips []*ufspb.IP) {
	if ips != nil && len(ips) > 0 {
		for _, v := range ips {
			logging.Debugf(ctx, "\tIP %s", v.GetId())
		}
		return
	}
	logging.Debugf(ctx, "\tNot found")
}

// ReplaceVlan replaces an old Vlan with new Vlan in datastore
//
// It does a delete of old vlan and create of new Vlan.
// All the steps are in done in a transaction to preserve consistency on failure.
// Before deleting the old Vlan, it will get all the resources referencing
// the old Vlan. It will update all the resources which were referencing
// the old Vlan(got in the last step) with new Vlan.
// Deletes the old Vlan.
// Creates the new Vlan.
// This will preserve data integrity in the system.
func ReplaceVlan(ctx context.Context, oldVlan *ufspb.Vlan, newVlan *ufspb.Vlan) (*ufspb.Vlan, error) {
	// TODO(eshwarn) : implement replace after user testing the tool
	return nil, nil
}

// validateCreateVlan validates if a vlan can be created
func validateCreateVlan(ctx context.Context, vlan *ufspb.Vlan) error {
	if err := resourceAlreadyExists(ctx, []*Resource{GetVlanResource(vlan.Name)}, nil); err != nil {
		return err
	}
	cidrBlock := vlan.GetVlanAddress()
	if cidrBlock != "" {
		vlans, err := configuration.QueryVlanByPropertyName(ctx, "cidr_block", cidrBlock, true)
		if err != nil {
			return err
		}
		if len(vlans) > 0 {
			return status.Errorf(codes.InvalidArgument, "cidr block %s is already occupied by %s", cidrBlock, vlans[0].GetName())
		}
	}
	return nil
}

// validateDeleteVlan validates if a Vlan can be deleted
//
// Checks if this Vlan(VlanID) is not referenced by other resources in the datastore.
// If there are any other references, delete will be rejected and an error will be returned.
func validateDeleteVlan(ctx context.Context, id string) error {
	machinelses, err := inventory.QueryMachineLSEByPropertyName(ctx, "vlan_id", id, true)
	if err != nil {
		return err
	}
	if len(machinelses) > 0 {
		return status.Errorf(codes.FailedPrecondition, "vlan %s is occupied by %d hosts, e.g. %#v", id, len(machinelses), ufsAPI.ParseResources(machinelses, "Name"))
	}
	vms, err := inventory.QueryVMByPropertyName(ctx, "vlan_id", id, true)
	if err != nil {
		return err
	}
	if len(vms) > 0 {
		return status.Errorf(codes.FailedPrecondition, "vlan %s is occupied by %d vms, e.g. %#v", id, len(vms), ufsAPI.ParseResources(vms, "Name"))
	}
	return nil
}

// validateUpdateVlan validates if a vlan can be updated
func validateUpdateVlan(ctx context.Context, vlan *ufspb.Vlan, mask *field_mask.FieldMask) error {
	if err := ResourceExist(ctx, []*Resource{GetVlanResource(vlan.Name)}, nil); err != nil {
		return err
	}
	return validateVlanUpdateMask(ctx, vlan, mask)
}

// validateVlanUpdateMask validates the update mask for vlan update
func validateVlanUpdateMask(ctx context.Context, vlan *ufspb.Vlan, mask *field_mask.FieldMask) error {
	if mask != nil {
		// validate the give field mask
		for _, path := range mask.Paths {
			switch path {
			case "name":
				return status.Error(codes.InvalidArgument, "validateVlanUpdateMask - name cannot be updated, delete and create a new vlan instead")
			case "update_time":
				return status.Error(codes.InvalidArgument, "validateVlanUpdateMask - update_time cannot be updated, it is a Output only field")
			case "description":
			case "resourceState":
			case "cidr_block":
				return status.Error(codes.InvalidArgument, "validateVlanUpdateMask - cidr_block cannot be updated, delete and create a new vlan instead")
			case "tags":
				// valid fields, nothing to validate.
			case "reserved_ips":
			default:
				return status.Errorf(codes.InvalidArgument, "validateVlanUpdateMask - unsupported update mask path %q", path)
			}
		}
	}
	return nil
}

// processVlanUpdateMask process the update field masks to do partial update
//
// Return a complete vlan object with updated and existing fields
func processVlanUpdateMask(oldVlan *ufspb.Vlan, vlan *ufspb.Vlan, mask *field_mask.FieldMask) (*ufspb.Vlan, error) {
	// update the fields in the existing/old vlan
	for _, path := range mask.Paths {
		switch path {
		case "description":
			oldVlan.Description = vlan.GetDescription()
		case "resourceState":
			oldVlan.ResourceState = vlan.GetResourceState()
		case "reserved_ips":
			oldVlan.ReservedIps = mergeIPs(oldVlan.ReservedIps, vlan.GetReservedIps())
		}
	}
	return oldVlan, nil
}

func getVlanHistoryClient(m *ufspb.Vlan) *HistoryClient {
	return &HistoryClient{
		stUdt: &stateUpdater{
			ResourceName: util.AddPrefix(util.VlanCollection, m.Name),
		},
		netUdt: &networkUpdater{
			Hostname: m.Name,
		},
	}
}

func updateIPTable(ctx context.Context, vlanName string, oldIPs, newIPs []string) error {
	oldIPMap := make(map[string]bool, len(oldIPs))
	for _, ip := range oldIPs {
		oldIPMap[ip] = true
	}
	newIPMap := make(map[string]bool, len(newIPs))
	for _, ip := range newIPs {
		newIPMap[ip] = true
	}
	// Get ips from reserved => non-reserved
	addBackIPs := make([]*ufspb.IP, 0)
	for _, ip := range oldIPs {
		if !newIPMap[ip] {
			addBackIPs = append(addBackIPs, util.FormatIP(vlanName, ip, false))
		}
	}
	// Get ips from non-reserved => reserved
	toRemoveIPs := make([]string, 0)
	for _, ip := range newIPs {
		if !oldIPMap[ip] {
			toRemoveIPs = append(toRemoveIPs, util.FormatIP(vlanName, ip, false).GetId())
		}
	}
	if _, err := configuration.BatchUpdateIPs(ctx, addBackIPs); err != nil {
		return err
	}
	if err := configuration.BatchDeleteIPs(ctx, toRemoveIPs); err != nil {
		return err
	}
	return nil
}
