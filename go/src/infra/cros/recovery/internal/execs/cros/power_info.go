// Copyright 2021 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cros

import (
	"context"
	"strconv"
	"strings"

	"go.chromium.org/luci/common/errors"
	"infra/cros/recovery/internal/execs"
)

// powerSupplyInfo holds info from power_supply_info.
type powerSupplyInfo struct {
	// The map of power_supply_info, e.g.,
	// {
	// 'Line Power':
	//  {
	//	 'online': 'yes',
	//	 'type': 'main'
	//  },
	// 'Battery':
	//  {
	// 	 'vendor': 'xyz',
	//	 'percentage': '100'
	//  }
	// }
	powerInfo map[string]map[string]string
}

// ReadPowerInfo initialize and return a new powerSupplyInfo sturct.
// Output of power_supply_info shows two devices, Line Power and Battery, with details of
// each device listed. This function parses the output into a dictionary,
// with key being the device name, and value being a dictionary of details of the device info.
//     Device: Line Power
//       online:                  no
//       type:                    Mains
//       voltage (V):             0
//       current (A):             0
//     Device: Battery
//       state:                   Discharging
//       percentage:              95.9276
//       technology:              Li-ion
func ReadPowerInfo(ctx context.Context, args *execs.RunArgs) (*powerSupplyInfo, error) {
	r := args.Access.Run(ctx, args.ResourceName, "power_supply_info")
	if r.ExitCode != 0 {
		return nil, errors.Reason("read power information: failed with code %d and %q", r.ExitCode, r.Stderr).Err()
	}
	return &powerSupplyInfo{
		powerInfo: getPowerSupplyInfoInMap(strings.TrimSpace(r.Stdout)),
	}, nil
}

// IsACOnline confirms the DUT is powered by AC.
func (p *powerSupplyInfo) IsACOnline() (bool, error) {
	if linePower, ok := p.powerInfo["Line Power"]; ok {
		if isOnline, ok := linePower["online"]; ok {
			return strings.ToLower(isOnline) == "yes", nil
		}
		return false, errors.Reason("ac online: no ac's online info found").Err()
	}
	return false, errors.Reason("ac online: no ac info found").Err()
}

// HasBattery confirms the DUT has a battery.
func (p *powerSupplyInfo) HasBattery() (bool, error) {
	if _, ok := p.powerInfo["Battery"]; ok {
		return ok, nil
	}
	return false, errors.Reason("has battery: no found").Err()
}

// IsBatteryDischarging confirms the DUT's battery is discharging.
func (p *powerSupplyInfo) IsBatteryDischarging() (bool, error) {
	if battery, ok := p.powerInfo["Battery"]; ok {
		if charging_state, ok := battery["state"]; ok {
			return charging_state == "Discharging", nil
		}
		return false, errors.Reason("battery discharging: no battery's state info found").Err()
	}
	return false, errors.Reason("battery discharging: no battery info found").Err()
}

// BatteryLevel returns the DUT's battery battery level.
func (p *powerSupplyInfo) BatteryLevel() (float64, error) {
	if battery, ok := p.powerInfo["Battery"]; ok {
		if percentage, ok := battery["percentage"]; ok {
			if batteryLevel, err := strconv.ParseFloat(percentage, 64); err != nil {
				return -1, errors.Annotate(err, "battery level").Err()
			} else {
				return batteryLevel, nil
			}
		}
		return -1, errors.Reason("battery level: no battery's percentage info found").Err()
	}
	return -1, errors.Reason("battery level: no battery").Err()
}

// getPowerSupplyInfoInMap is a helper function to get power supply information for ReadPowerInfo().
func getPowerSupplyInfoInMap(rawOutput string) map[string]map[string]string {
	info := make(map[string]map[string]string)
	device_name := ""
	var device_info map[string]string
	temp_result := strings.Split(rawOutput, "\n")
	for _, eachLine := range temp_result {
		pairs := strings.Split(eachLine, ":")
		if len(pairs) != 2 {
			continue
		}
		key := strings.TrimSpace(pairs[0])
		val := strings.TrimSpace(pairs[1])
		if key == "Device" {
			if device_name != "" {
				info[device_name] = device_info
			}
			device_name = val
			device_info = make(map[string]string)
		} else if device_info != nil {
			device_info[key] = val
		}
	}
	if _, ok := info[device_name]; !ok && device_name != "" {
		info[device_name] = device_info
	}
	return info
}