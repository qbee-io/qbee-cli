package client

import (
	"context"
	"fmt"
	"net/http"
)

type DeviceStatus struct {
	RemoteAccess bool   `json:"remote_access"`
	Edge         string `json:"edge,omitempty"`
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
