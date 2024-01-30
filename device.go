package client

import (
	"context"
	"fmt"
	"net/http"
)

// EdgeVersion indicates the version of the edge server that the device is connected to.
type EdgeVersion uint8

const (
	EdgeVersionOpenVPN = 0
	EdgeVersionNative  = 1
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
