package client

import (
	"context"
	"fmt"
	"net/http"
)

// PendingDevice represents a device that is pending approval.
type PendingDevice struct {
	// PublicKeyDigest - device's public key digest (used to identify devices).
	PublicKeyDigest string `json:"pub_key_digest"`

	// Approved is true when the device is approved.
	Approved bool `json:"approved"`

	// ApprovedTimestamp in UnixNano.
	ApprovedTimestamp int64 `json:"approved_timestamp,omitempty"`

	// Certificate - base64-encoded (PEM) signed device certificate.
	Certificate string `json:"cert,omitempty"`

	// GroupID where device is assigned.
	GroupID string `json:"group_id,omitempty"`

	// RemoteAddress from which device was bootstrapped.
	RemoteAddress string `json:"remoteaddress"`

	// Device creation timestamp (in seconds).
	Timestamp int64 `json:"timestamp"`

	// Host - The name of the current host, according to the kernel.
	// It is undefined whether this is qualified or unqualified with a domain name.
	Host string `json:"host"`

	// FQHost - The fully qualified name of the host (e.g. "device1.example.com").
	FQHost string `json:"fqhost"`

	// UQHost - The unqualified name of the host (e.g. "device1").
	UQHost string `json:"uqhost"`

	// DeviceName - The custom name for the device.
	DeviceName string `json:"device_name,omitempty"`

	// HardwareMAC - This contains the MAC address of the named interface map[interface]macAddress.
	// Note: The keys in this array are canonified.
	// For example, the entry for wlan0.1 would be found under the wlan0_1 key.
	//
	// Example:
	// {
	// 	"ens1": "52:54:00:4a:db:ee",
	//  "qbee0": "00:00:00:00:00:00"
	// }
	HardwareMAC map[string]string `json:"hardware_mac"`

	// IPDefault - All four octets of the IPv4 address of the first system interface.
	//Note: If the system has a single ethernet interface, this variable will contain the IPv4 address.
	// However, if the system has multiple interfaces, then this variable will simply be the IPv4 address of the first
	// interface in the list that has an assigned address.
	// Use IPv4[interface_name] for details on obtaining the IPv4 addresses of all interfaces on a system.
	IPDefault string `json:"ip_default"`

	// IPv4 - All IPv4 addresses of the system mapped by interface name.
	// Example:
	// {
	//	"ens1": "192.168.122.239",
	//	"qbee0": "100.64.39.78"
	// }
	IPv4 map[string]string `json:"ipv4"`

	// RawPublicKey of the device as slice of PEM-encoded lines.
	// Example:
	// []string{
	//    "-----BEGIN PUBLIC KEY-----",
	//    "MIGbMBAGByqGSM49AgEGBSuBBAAjA4GGAAQBvDALiaU+dyvd1DhMUCEXnuX4h5r5",
	//    "ikBVNSl88QBtBoxtQy1XsCJ7Dm/tzoQ1YPYT80oVTdExK/oFnZFvI89SX8sBN89L",
	//    "Y8q+8BBQPLf1nA3DG7apq7xq11Zkpde2eK0pCUG7nZPisXlU96C44NLE62TzDYEZ",
	//    "RYkhJQhFeNOlFSpF/xA=",
	//    "-----END PUBLIC KEY-----"
	// }
	RawPublicKey []string `json:"pub_key,omitempty"`
}

const pendingHostPath = "/api/v2/pendinghost"

// GetPendingDevices returns the pending devices list.
func (cli *Client) GetPendingDevices(ctx context.Context, pendingOnly bool) ([]PendingDevice, error) {
	var devices []PendingDevice

	path := pendingHostPath
	if pendingOnly {
		path += "?pendingonly=1"
	}

	if err := cli.Call(ctx, http.MethodGet, path, nil, &devices); err != nil {
		return nil, err
	}

	return devices, nil
}

// ApprovePendingDeviceRequest represents a request to approve a pending device.
type ApprovePendingDeviceRequest struct {
	// NodeID is the node ID of the device to approve.
	NodeID string `json:"node_id"`
}

// ApprovePendingDevice approves a pending device with provided node ID.
func (cli *Client) ApprovePendingDevice(ctx context.Context, nodeID string) error {
	data := ApprovePendingDeviceRequest{
		NodeID: nodeID,
	}

	if err := cli.Call(ctx, http.MethodPost, pendingHostPath, data, nil); err != nil {
		return err
	}

	return nil
}

// RejectPendingDevice rejects a pending device with provided node ID.
func (cli *Client) RejectPendingDevice(ctx context.Context, nodeID string) error {
	path := pendingHostPath + "/" + nodeID

	if err := cli.Call(ctx, http.MethodDelete, path, nil, nil); err != nil {
		return err
	}

	return nil
}

const removeApprovedHostPath = "/api/v2/removeapprovedhost/%s"

// RemoveApprovedHost rejects a pending device with provided node ID.
func (cli *Client) RemoveApprovedHost(ctx context.Context, nodeID string) error {
	path := fmt.Sprintf(removeApprovedHostPath, nodeID)

	if err := cli.Call(ctx, http.MethodDelete, path, nil, nil); err != nil {
		return err
	}

	return nil
}
