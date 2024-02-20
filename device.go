// Copyright 2024 qbee.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"fmt"
	"net/http"
)

// EdgeVersion indicates the version of the edge server that the device is connected to.
type EdgeVersion uint8

const (
	// EdgeVersionOpenVPN indicates that the device is connected to an OpenVPN edge.
	EdgeVersionOpenVPN = 0

	// EdgeVersionNative indicates that the device is connected to a native qbee remote access edge.
	EdgeVersionNative = 1
)

// DeviceStatus is the status of a device.
type DeviceStatus struct {
	// UUID is the UUID of the device.
	UUID string `json:"uuid"`

	// RemoteAccess is true if the device is connected to the edge.
	RemoteAccess bool `json:"remote_access"`

	// Edge is the edge host that the device is connected to.
	// This field is only set if RemoteAccess is true.
	// Format is <edge-host>:<edge-port>/edge/<edge-id>
	Edge string `json:"edge,omitempty"`

	// EdgeVersion is the version of the edge that the device is connected to.
	// This field is only set if RemoteAccess is true.
	// 0 - for OpenVPN edge
	// 1 - for native qbee remote access
	EdgeVersion EdgeVersion `json:"edge_version,omitempty"`
}

// GetDeviceStatus returns device status.
func (cli *Client) GetDeviceStatus(ctx context.Context, deviceID string) (*DeviceStatus, error) {
	deviceStatus := new(DeviceStatus)

	path := fmt.Sprintf("/api/v2/device/%s/status", deviceID)

	if err := cli.Call(ctx, http.MethodGet, path, nil, deviceStatus); err != nil {
		return nil, err
	}

	return deviceStatus, nil
}
