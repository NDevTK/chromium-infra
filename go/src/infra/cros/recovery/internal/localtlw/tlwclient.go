// Copyright 2021 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package localtlw provides local implementation of TLW Access.
package localtlw

import (
	"context"
	"fmt"

	"go.chromium.org/chromiumos/config/go/api/test/xmlrpc"
	"go.chromium.org/luci/common/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fleet "infra/appengine/crosskylabadmin/api/fleet/v1"
	"infra/cros/recovery/internal/localtlw/dutinfo"

	"infra/cros/dutstate"
	tlwio "infra/cros/recovery/internal/localtlw/io"
	"infra/cros/recovery/internal/localtlw/localproxy"
	"infra/cros/recovery/internal/localtlw/rpm"
	"infra/cros/recovery/internal/localtlw/servod"
	"infra/cros/recovery/internal/localtlw/ssh"
	"infra/cros/recovery/internal/log"
	"infra/cros/recovery/tlw"
	"infra/libs/sshpool"
	ufspb "infra/unifiedfleet/api/v1/models"
	ufslab "infra/unifiedfleet/api/v1/models/chromeos/lab"
	ufsAPI "infra/unifiedfleet/api/v1/rpc"
	ufsUtil "infra/unifiedfleet/app/util"
)

const (
	// gsCrosImageBucket is the base URL for the Google Storage bucket for
	// ChromeOS image archives.
	gsCrosImageBucket = "gs://chromeos-image-archive"
	// tlwPort is default port used to run TLW on the drones.
	tlwPort = 7151
	// tlsPort is default port used to run TLS on the drones.
	tlsPort = 7152
)

// UFSClient is a client that knows how to work with UFS RPC methods.
type UFSClient interface {
	// GetSchedulingUnit retrieves the details of the SchedulingUnit.
	GetSchedulingUnit(ctx context.Context, req *ufsAPI.GetSchedulingUnitRequest, opts ...grpc.CallOption) (rsp *ufspb.SchedulingUnit, err error)
	// GetChromeOSDeviceData retrieves requested Chrome OS device data from the UFS and inventoryV2.
	GetChromeOSDeviceData(ctx context.Context, req *ufsAPI.GetChromeOSDeviceDataRequest, opts ...grpc.CallOption) (rsp *ufspb.ChromeOSDeviceData, err error)
	// UpdateDutState updates the state config for a DUT
	UpdateDutState(ctx context.Context, in *ufsAPI.UpdateDutStateRequest, opts ...grpc.CallOption) (*ufslab.DutState, error)
}

// CSAClient is a client that knows how to respond to the GetStableVersion RPC call.
type CSAClient interface {
	GetStableVersion(ctx context.Context, in *fleet.GetStableVersionRequest, opts ...grpc.CallOption) (*fleet.GetStableVersionResponse, error)
}

// tlwClient holds data and represents the local implementation of TLW Access interface.
type tlwClient struct {
	csaClient  CSAClient
	ufsClient  UFSClient
	sshPool    *sshpool.Pool
	servodPool *servod.Pool
	// Cache received devices from inventory
	devices map[string]*ufspb.ChromeOSDeviceData
}

// New build new local TLW Access instance.
func New(ufs UFSClient, csac CSAClient) (tlw.Access, error) {
	c := &tlwClient{
		ufsClient:  ufs,
		csaClient:  csac,
		sshPool:    sshpool.New(ssh.SSHConfig()),
		servodPool: servod.NewPool(),
		devices:    make(map[string]*ufspb.ChromeOSDeviceData),
	}
	return c, nil
}

// Close closes all used resources.
func (c *tlwClient) Close() error {
	if err := c.sshPool.Close(); err != nil {
		return errors.Annotate(err, "tlw client").Err()
	}
	return c.servodPool.Close()
}

// Ping performs ping by resource name.
func (c *tlwClient) Ping(ctx context.Context, resourceName string, count int) error {
	return ping(resourceName, count)
}

// Run executes command on device by SSH related to resource name.
func (c *tlwClient) Run(ctx context.Context, resourceName, command string) *tlw.RunResult {
	return ssh.Run(ctx, c.sshPool, localproxy.BuildAddr(resourceName), command)
}

// InitServod initiates servod daemon on servo-host.
func (c *tlwClient) InitServod(ctx context.Context, req *tlw.InitServodRequest) error {
	dd, err := c.getDevice(ctx, req.Resource)
	if err != nil {
		return errors.Annotate(err, "init servod %q", req.Resource).Err()
	}
	dut := dd.GetLabConfig().GetChromeosMachineLse().GetDeviceLse().GetDut()
	if dut == nil {
		return errors.Reason("init servod %q: dut is not found", req.Resource).Err()
	}
	servo := dut.GetPeripherals().GetServo()
	if servo == nil {
		return errors.Reason("init servod %q: servo is not found", req.Resource).Err()
	}
	s, err := c.servodPool.Get(
		localproxy.BuildAddr(servo.GetServoHostname()),
		servo.GetServoPort(),
		func() ([]string, error) {
			return dutinfo.GenerateServodParams(dd, req.Options)
		})
	if err != nil {
		return errors.Annotate(err, "init servod %q", req.Resource).Err()
	}
	if err := s.Prepare(ctx, c.sshPool); err != nil {
		return errors.Annotate(err, "init servod %q", req.Resource).Err()
	}
	return nil
}

// StopServod stops servod daemon on servo-host.
func (c *tlwClient) StopServod(ctx context.Context, resourceName string) error {
	dd, err := c.getDevice(ctx, resourceName)
	if err != nil {
		return errors.Annotate(err, "stop servod %q", resourceName).Err()
	}
	dut := dd.GetLabConfig().GetChromeosMachineLse().GetDeviceLse().GetDut()
	if dut == nil {
		return errors.Reason("stop servod %q: dut is not found", resourceName).Err()
	}
	servo := dut.GetPeripherals().GetServo()
	if servo == nil {
		log.Debug(ctx, "Stop servod %q: servo is not specified", resourceName)
		return nil
	}
	s, err := c.servodPool.Get(localproxy.BuildAddr(servo.GetServoHostname()), servo.GetServoPort(), nil)
	if err != nil {
		return errors.Annotate(err, "stop servod %q", resourceName).Err()
	}
	if err := s.Stop(ctx, c.sshPool); err != nil {
		return errors.Annotate(err, "stop servod %q", resourceName).Err()
	}
	return nil
}

// CallServod executes a command on servod related to resource name.
// Commands will be run against servod on servo-host.
func (c *tlwClient) CallServod(ctx context.Context, req *tlw.CallServodRequest) *tlw.CallServodResponse {
	// Translator to convert error to response structure.
	fail := func(err error) *tlw.CallServodResponse {
		return &tlw.CallServodResponse{
			Value: &xmlrpc.Value{
				ScalarOneof: &xmlrpc.Value_String_{
					String_: fmt.Sprintf("call servod %q: %s", req.Resource, err),
				},
			},
			Fault: true,
		}
	}
	dd, err := c.getDevice(ctx, req.Resource)
	if err != nil {
		return fail(err)
	}
	dut := dd.GetLabConfig().GetChromeosMachineLse().GetDeviceLse().GetDut()
	if dut == nil {
		log.Debug(ctx, "Call servod %q: dut is not found", req.Resource)
		return nil
	}
	servo := dut.GetPeripherals().GetServo()
	if servo == nil {
		log.Debug(ctx, "Call servod %q: servo is not specified", req.Resource)
		return nil
	}
	s, err := c.servodPool.Get(
		localproxy.BuildAddr(servo.GetServoHostname()),
		servo.GetServoPort(), func() ([]string, error) {
			return dutinfo.GenerateServodParams(dd, req.Options)
		})
	if err != nil {
		return fail(err)
	}
	res, err := s.Call(ctx, c.sshPool, req)
	if err != nil {
		return fail(err)
	}
	return res
}

// CopyFileTo copies file to remote device from local.
func (c *tlwClient) CopyFileTo(ctx context.Context, req *tlw.CopyRequest) error {
	if err := tlwio.CopyFileTo(ctx, c.sshPool, req); err != nil {
		return errors.Annotate(err, "copy file to").Err()
	}
	return nil
}

// CopyFileFrom copies file from remote device to local.
func (c *tlwClient) CopyFileFrom(ctx context.Context, req *tlw.CopyRequest) error {
	if err := tlwio.CopyFileFrom(ctx, c.sshPool, req); err != nil {
		return errors.Annotate(err, "copy file from").Err()
	}
	return nil
}

// CopyDirectoryTo copies directory to remote device from local, recursively.
func (c *tlwClient) CopyDirectoryTo(ctx context.Context, req *tlw.CopyRequest) error {
	if err := tlwio.CopyDirectoryTo(ctx, c.sshPool, req); err != nil {
		return errors.Annotate(err, "copy directory to").Err()
	}
	return nil
}

// CopyDirectoryFrom copies directory from remote device to local, recursively.
func (c *tlwClient) CopyDirectoryFrom(ctx context.Context, req *tlw.CopyRequest) error {
	if err := tlwio.CopyDirectoryFrom(ctx, c.sshPool, req); err != nil {
		return errors.Annotate(err, "copy directory from").Err()
	}
	return nil
}

// SetPowerSupply manages power supply for requested.
func (c *tlwClient) SetPowerSupply(ctx context.Context, req *tlw.SetPowerSupplyRequest) *tlw.SetPowerSupplyResponse {
	if req == nil || req.Resource == "" {
		return &tlw.SetPowerSupplyResponse{
			Status: tlw.PowerSupplyResponseStatusError,
			Reason: "resource is not specified",
		}
	}
	dd, err := c.getDevice(ctx, req.Resource)
	if err != nil {
		return &tlw.SetPowerSupplyResponse{
			Status: tlw.PowerSupplyResponseStatusError,
			Reason: err.Error(),
		}
	}
	hostname, outlet := dutinfo.GetRpmInfo(dd)
	log.Debug(ctx, "Set power supply %s: has rpm info %s:%s.", req.Resource, hostname, outlet)
	if hostname == "" || outlet == "" {
		return &tlw.SetPowerSupplyResponse{
			Status: tlw.PowerSupplyResponseStatusNoConfig,
			Reason: "power supply config missing or incorrect",
		}
	}
	var s rpm.PowerState
	switch req.State {
	case tlw.PowerSupplyActionOn:
		s = rpm.PowerStateOn
	case tlw.PowerSupplyActionOff:
		s = rpm.PowerStateOff
	case tlw.PowerSupplyActionCycle:
		s = rpm.PowerStateCycle
	default:
		return &tlw.SetPowerSupplyResponse{
			Status: tlw.PowerSupplyResponseStatusError,
			Reason: fmt.Sprintf("unknown rpm state: %s", string(req.State)),
		}
	}
	log.Debug(ctx, "Set power supply %s: state: %q for %s:%s.", req.Resource, s, hostname, outlet)
	rpmReq := &rpm.RPMPowerRequest{
		Hostname:          dd.GetLabConfig().GetName(),
		PowerUnitHostname: hostname,
		PowerunitOutlet:   outlet,
		State:             s,
	}
	if err := rpm.SetPowerState(ctx, rpmReq); err != nil {
		return &tlw.SetPowerSupplyResponse{
			Status: tlw.PowerSupplyResponseStatusError,
			Reason: err.Error(),
		}
	}
	return &tlw.SetPowerSupplyResponse{
		Status: tlw.PowerSupplyResponseStatusOK,
	}
}

// GetCacheUrl provides URL to download requested path to file.
// URL will use to download image to USB-drive and provisioning.
func (c *tlwClient) GetCacheUrl(ctx context.Context, resourceName, filePath string) (string, error) {
	// TODO(otabek@): Add logic to understand local file and just return it back.
	addr := fmt.Sprintf("0.0.0.0:%d", tlwPort)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return "", errors.Annotate(err, "connect to background TLW").Err()
	}
	defer func() { conn.Close() }()
	return CacheForDut(ctx, conn, filePath, resourceName)
}

// ListResourcesForUnit provides list of resources names related to target unit.
func (c *tlwClient) ListResourcesForUnit(ctx context.Context, name string) ([]string, error) {
	if name == "" {
		return nil, errors.Reason("list resources: unit name is expected").Err()
	}
	dd, err := c.ufsClient.GetChromeOSDeviceData(ctx, &ufsAPI.GetChromeOSDeviceDataRequest{
		Hostname: name,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			log.Debug(ctx, "List resources %q: record not found.", name)
		} else {
			return nil, errors.Reason("list resources %q", name).Err()
		}
	} else if dd.GetLabConfig() == nil {
		return nil, errors.Reason("list resources %q: device data is empty", name).Err()
	} else {
		log.Debug(ctx, "List resources %q: cached received device.", name)
		c.devices[name] = dd
		return []string{dd.GetLabConfig().GetName()}, nil
	}
	suName := ufsUtil.AddPrefix(ufsUtil.SchedulingUnitCollection, name)
	log.Debug(ctx, "list resources %q: trying to find scheduling unit by name %q.", name, suName)
	su, err := c.ufsClient.GetSchedulingUnit(ctx, &ufsAPI.GetSchedulingUnitRequest{
		Name: suName,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, errors.Annotate(err, "list resources %q: record not found", name).Err()
		}
		return nil, errors.Annotate(err, "list resources %q", name).Err()
	}
	var resourceNames []string
	for _, hostname := range su.GetMachineLSEs() {
		resourceNames = append(resourceNames, hostname)
	}
	return resourceNames, nil
}

// GetDut provides DUT info per requested resource name from inventory.
func (c *tlwClient) GetDut(ctx context.Context, name string) (*tlw.Dut, error) {
	dd, err := c.getDevice(ctx, name)
	if err != nil {
		return nil, errors.Annotate(err, "get DUT %q", name).Err()
	}
	dut, err := dutinfo.ConvertDut(dd)
	if err != nil {
		return nil, errors.Annotate(err, "get DUT %q", name).Err()
	}
	gv, err := c.getStableVersion(ctx, dut)
	if err != nil {
		log.Info(ctx, "Get DUT %q: failed to receive stable-version. Error: %s", name, err)
	} else {
		dut.StableVersion = gv
	}
	return dut, nil
}

// getDevice receives device from inventory.
func (c *tlwClient) getDevice(ctx context.Context, name string) (*ufspb.ChromeOSDeviceData, error) {
	if d, ok := c.devices[name]; ok {
		log.Debug(ctx, "Get device %q: received from cache.", name)
		return d, nil
	}
	req := &ufsAPI.GetChromeOSDeviceDataRequest{Hostname: name}
	dd, err := c.ufsClient.GetChromeOSDeviceData(ctx, req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, errors.Reason("get device %q: record not found", name).Err()
		}
		return nil, errors.Annotate(err, "get device %q", name).Err()
	} else if dd.GetLabConfig() == nil {
		return nil, errors.Reason("get device %q: received empty data", name).Err()
	}
	c.devices[name] = dd
	log.Debug(ctx, "Get device %q: cached received device.", name)
	return dd, nil
}

// getStableVersion receives stable versions of device.
func (c *tlwClient) getStableVersion(ctx context.Context, dut *tlw.Dut) (*tlw.StableVersion, error) {
	req := &fleet.GetStableVersionRequest{Hostname: dut.Name}
	res, err := c.csaClient.GetStableVersion(ctx, req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, errors.Reason("get stable-version %q: record not found", dut.Name).Err()
		}
		return nil, errors.Annotate(err, "get stable-version %q", dut.Name).Err()
	}
	if res.GetCrosVersion() == "" {
		return nil, errors.Reason("get stable-version %q: version is empty", dut.Name).Err()
	}
	return &tlw.StableVersion{
		CrosImage:           fmt.Sprintf("%s-release/%s", dut.Board, res.GetCrosVersion()),
		CrosFirmwareVersion: res.GetFirmwareVersion(),
		CrosFirmwareImage:   res.GetFaftVersion(),
	}, nil
}

// UpdateDut updates DUT info into inventory.
func (c *tlwClient) UpdateDut(ctx context.Context, dut *tlw.Dut) error {
	if dut == nil {
		return errors.Reason("update DUT: DUT is not provided").Err()
	}
	dd, err := c.getDevice(ctx, dut.Name)
	if err != nil {
		return errors.Annotate(err, "update DUT %q", dut.Name).Err()
	}
	dutID := dd.GetMachine().GetName()
	req, err := dutinfo.CreateUpdateDutRequest(dutID, dut)
	if err != nil {
		return errors.Annotate(err, "update DUT %q", dut.Name).Err()
	}
	log.Debug(ctx, "Update DUT: update request: %s", req)
	if _, err := c.ufsClient.UpdateDutState(ctx, req); err != nil {
		return errors.Annotate(err, "update DUT %q", dut.Name).Err()
	}
	delete(c.devices, dut.Name)
	if ufs, ok := c.ufsClient.(dutstate.UFSClient); ok {
		if err := dutstate.Update(ctx, ufs, dut.Name, dut.State); err != nil {
			return errors.Annotate(err, "update DUT %q", dut.Name).Err()
		}
	} else {
		return errors.Reason("update DUT %q: dutstate.UFSClient interface is not implemented by client", dut.Name).Err()
	}
	return nil
}

// Provision triggers provisioning of the device.
func (c *tlwClient) Provision(ctx context.Context, req *tlw.ProvisionRequest) error {
	if req == nil {
		return errors.Reason("provision: request is empty").Err()
	}
	if req.Resource == "" {
		return errors.Reason("provision: resource is not specified").Err()
	}
	if req.SystemImagePath == "" {
		return errors.Reason("provision: system image path is not specified").Err()
	}
	log.Debug(ctx, "Started provisioning by TLS: %s", req)
	addr := fmt.Sprintf("0.0.0.0:%d", tlsPort)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return errors.Annotate(err, "provision: connect to TLS").Err()
	}
	defer func() { conn.Close() }()
	err = TLSProvision(ctx, conn, req)
	return errors.Annotate(err, "provision").Err()
}
