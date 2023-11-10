// Copyright 2023 qbee.io
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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/qbee-io/qbee-cli/client/config"
)

const deviceInventoryPath = "/api/v2/inventory"

type SystemInfo struct {
	// Class - This variable contains the name of the hard-class category for this host,
	// (i.e. its top level operating system type classification, e.g. "linux").
	Class string `json:"class"`

	// OS - The name of the operating system according to the kernel (e.g. "linux").
	OS string `json:"os"`

	// OSType - Another name for the operating system (e.g. "linux_x86_64").
	OSType string `json:"ostype"`

	// Version - The version of the running kernel. On Linux, this corresponds to the output of uname -v.
	// Example: "#58-Ubuntu SMP Thu Oct 13 08:03:55 UTC 2022".
	Version string `json:"version"`

	// Architecture - The variable gives the kernel's short architecture description (e.g. "x86_64").
	Architecture string `json:"arch"`

	// LongArchitecture - The long architecture name for this system kernel.
	// This name is sometimes quite unwieldy but can be useful for logging purposes.
	// Example: "linux_x86_64_5_15_0_52_generic__58_Ubuntu_SMP_Thu_Oct_13_08_03_55_UTC_2022"
	LongArchitecture string `json:"long_arch"`

	// Release - The kernel release of the operating system (e.g. "5.15.0-52-generic").
	Release string `json:"release"`

	// Flavor - A variable containing an operating system identification string that is used to determine
	// the current release of the operating system in a form that can be used as a label in naming.
	// This is used, for instance, to detect which package name to choose when updating software binaries.
	// Example: "ubuntu_22"
	Flavor string `json:"flavor"`

	// BootTime represents system boot time (as Unix timestamp string, e.g. "1586144402")
	BootTime string `json:"boot_time"`

	// CPUs - A variable containing the number of CPU cores detected. On systems which provide virtual cores,
	// it is set to the total number of virtual, not physical, cores.
	// In addition, on a single-core system the class 1_cpu is set, and on multicore systems the class n_cpus is set,
	// where n is the number of cores identified (e.g. "4").
	CPUs string `json:"cpus"`

	// CpuSerialNumber - the serial number of the CPU (e.g. "0000000000000000").
	CPUSerialNumber string `json:"cpu_sn"`

	// CPURevision - the revision of the CPU (e.g. "10")
	CPURevision string `json:"cpu_rev"`

	// CPUHardware - the CPU hardware description (e.g. "Freescale i.MX6 Quad/DualLite (Device Tree)").
	CPUHardware string `json:"cpu_hw"`

	// Host - The name of the current host, according to the kernel.
	// It is undefined whether this is qualified or unqualified with a domain name.
	Host string `json:"host"`

	// FQHost - The fully qualified name of the host (e.g. "device1.example.com").
	FQHost string `json:"fqhost"`

	// UQHost - The unqualified name of the host (e.g. "device1").
	UQHost string `json:"uqhost"`

	// Interface - The assumed (default) name of the main system interface on this host.
	Interface string `json:"interface"`

	// HardwareMAC - This contains the MAC address of the named interface map[interface]macAddress.
	// For example, the entry for wlan0.1 would be found under the wlan0_1 key.
	//
	// Example:
	// {
	// 	"ens1": "52:54:00:4a:db:ee",
	//  "qbee0": "00:00:00:00:00:00"
	// }
	HardwareMAC map[string]string `json:"hardware_mac"`

	// InterfaceFlags - Contains a space separated list of the flags of the named interfaces.
	// The following device flags are supported:
	//    up
	//    broadcast
	//    debug
	//    loopback
	//    pointopoint
	//    notrailers
	//    running
	//    noarp
	//    promisc
	//    allmulti
	//    multicast
	//
	// Example:
	// {
	// 	"ens1": "up broadcast running multicast",
	//  "qbee0": "up pointopoint running noarp multicast"
	// }
	InterfaceFlags map[string]string `json:"interface_flags"`

	// IPAddresses - A system list of IP addresses currently in use by the system (e.g: "100.64.39.78").
	IPAddresses string `json:"ip_addresses"`

	// IPv4First - All four octets of the IPv4 address of the first system interface.
	//Note: If the system has a single ethernet interface, this variable will contain the IPv4 address.
	// However, if the system has multiple interfaces, then this variable will simply be the IPv4 address of the first
	// interface in the list that has an assigned address.
	// Use IPv4[interface_name] for details on obtaining the IPv4 addresses of all interfaces on a system.
	IPv4First string `json:"ipv4_first"`

	// IPv4 - All IPv4 addresses of the system mapped by interface name.
	// Example:
	// {
	//	"ens1": "192.168.122.239",
	//	"qbee0": "100.64.39.78"
	// }
	IPv4 map[string]string `json:"ipv4"`

	// RemoteAddress - remote client address from which the inventory was reported (e.g. "1.2.3.4").
	RemoteAddress string `json:"remoteaddr"`

	// LastConfigUpdate - unix timestamp of the last config update (e.g. "1586144402").
	LastConfigUpdate string `json:"last_config_update"`

	// LastConfigCommitID - last applied config commit SHA
	// (e.g. "6c07b6d021a015329b1815ec954cca6d8c4973c3b574202401dad448e8cdd0f5").
	LastConfigCommitID string `json:"last_config_commit_id"`

	// VPNIndex - defines numeric ID of the VPN server to which the device is connected.
	VPNIndex string `json:"vpn_idx"`

	// AgentVersion version of the agent which is currently running on the device.
	AgentVersion string `json:"cf_version"`
}

// DeviceInventory represents basic device inventory information.
type DeviceInventory struct {
	// ID is the device ID for which the inventory was reported.
	ID string `json:"pub_key_digest"`

	// SystemInfo contains system information.
	SystemInfo SystemInfo `json:"system"`

	// Settings contains device settings.
	// DEPRECATED: Use PushedConfig.BundleData.Settings instead.
	Settings config.Settings `json:"settings"`

	// ConfigPropagated is set to true if the device has received the latest configuration.
	ConfigPropagated bool `json:"config_propagated"`

	// LastReported is the unix timestamp when the device checked-in with the platform.
	LastReported int `json:"last_reported"`

	// HeartbeatExpirationSoft is the unix timestamp when the device will be marked as delayed.
	HeartbeatExpirationSoft int `json:"exp_soft"`

	// HeartbeatExpirationHard is the unix timestamp when the device will be marked as offline.
	HeartbeatExpirationHard int `json:"exp_hard"`

	// DeviceCommitSHA is the SHA of the commit that is currently running on the device.
	DeviceCommitSHA string `json:"device_commit_sha"`

	// Attributes contains currently defined device attributes.
	Attributes DeviceAttributes `json:"attributes"`

	// Tags of the device.
	// When set, tags are validated against the following regular expression: ^[a-z0-9:-]{3,64}$
	Tags []string `json:"tags"`

	// Status of the device.
	// Possible values: "online", "delayed", "offline"
	Status string `json:"status"`

	// PushedConfig contains the configuration that is expected to be on the device.
	// See ConfigPropagated to check if the device has already received it.
	PushedConfig config.Pushed `json:"pushed_config"`
}

// InventoryListSearch defines search parameters for InventoryListQuery.
type InventoryListSearch struct {
	// NodeID - device public key digest (hex-encoded SHA256)
	NodeID string `json:"node_id,omitempty"`

	// UUID - device UUID (legacy identifier used by the remote access subsystem)
	UUID string `json:"uuid,omitempty"`

	// Title - fqhost, remoteaddr or device title (can be set in the attributes)
	Title string `json:"title,omitempty"`

	// Flavor - A variable containing an operating system identification string that is used to determine
	// the current release of the operating system in a form that can be used as a label in naming.
	// This is used, for instance, to detect which package name to choose when updating software binaries.
	// Example: "ubuntu_22"
	Flavor string `json:"flavor,omitempty"`

	// Architecture - The variable containing the kernel's short architecture description (e.g. "x86_64")
	Architecture string `json:"arch,omitempty"`

	// IPAddress match on ipv4_first, ip_addresses, remoteaddr or ipv4 fields from inventory (exact match)
	IPAddress string `json:"ip,omitempty"`

	// MACAddress match on hardware_mac field from inventory (exact match)
	MACAddress string `json:"mac,omitempty"`

	// Description partial match on device description.
	Description string `json:"description,omitempty"`

	// Latitude and Longitude match on device location.
	Latitude  string `json:"latitude,omitempty"`
	Longitude string `json:"longitude,omitempty"`

	// Country match on device location.
	Country string `json:"country,omitempty"`

	// City match on device location.
	City string `json:"city,omitempty"`

	// Zip match on device location.
	Zip string `json:"zip,omitempty"`

	// Address match on device location.
	Address string `json:"address,omitempty"`

	// DeviceAttribute - match on device attributes:
	// - flavor,
	// - architecture,
	// - description,
	// - lat/long,
	// - country,
	// - city,
	// - zip,
	// - address,
	// - ipv4_first,
	// - ip_addresses,
	// - remoteaddr,
	// - ipv4
	// using OR condition and substring match
	DeviceAttribute string `json:"device_attribute,omitempty"`

	// Ancestors - array of parent nodes (exact match)
	Ancestors []string `json:"ancestor_ids,omitempty"`

	// Tags - array of tags which MUST be preset (exact match)
	Tags []string `json:"tags,omitempty"`

	// CommitSHA - match on active config commit sha (partial match)
	CommitSHA string `json:"commit_sha,omitempty"`
}

type InventoryReportType string

const (
	// InventoryReportTypeDetailed - detailed inventory report incl. complete device attributes and system info.
	InventoryReportTypeDetailed InventoryReportType = "detailed"

	// InventoryReportTypeShort - short inventory report, incl. only basic device information:
	// system info:
	// - fqhost
	// - remoteaddr
	// - ipv4
	// - last_config_commit_id
	// - last_policy_update
	// device attributes:
	// - device name
	InventoryReportTypeShort InventoryReportType = "short"
)

// InventoryListQuery defines query parameters for InventoryList.
type InventoryListQuery struct {
	Search InventoryListSearch

	// SortField defines field used to sort, 'title' by default.
	// Supported sort fields:
	// - title (fqhost or remoteaddr - depending on what is available)
	// - device_name (from attribute)
	// - last_reported_time
	// - last_policy_update
	// - commit_sha
	SortField string

	// SortDirection defines sort direction, 'desc' by default.
	SortDirection string

	// ItemsPerPage defines maximum number of records in result, default 30, max 1000
	ItemsPerPage int

	// Offset defines offset of the first record in result, default 0
	Offset int

	// ReportType defines format of the response payload.
	ReportType InventoryReportType
}

// String returns string representation of InventoryListQuery which can be used as query string.
func (q InventoryListQuery) String() (string, error) {
	values := make(url.Values)

	searchQueryJSON, err := json.Marshal(q.Search)
	if err != nil {
		return "", err
	}

	values.Set("search", string(searchQueryJSON))

	if q.SortField != "" {
		values.Set("sort_field", q.SortField)
	}

	if q.SortDirection != "" {
		values.Set("sort_direction", q.SortDirection)
	}

	if q.ItemsPerPage != 0 {
		values.Set("items_per_page", fmt.Sprintf("%d", q.ItemsPerPage))
	}

	if q.Offset != 0 {
		values.Set("offset", fmt.Sprintf("%d", q.Offset))
	}

	if q.ReportType != "" {
		values.Set("report_type", string(q.ReportType))
	}

	return values.Encode(), nil
}

// InventoryListItem contains device inventory information.
type InventoryListItem struct {
	// NodeID - device public key digest (hex-encoded SHA256)
	NodeID string `json:"node_id"`

	// PubKeyDigest - same as NodeID (for backwards compatibility)
	PubKeyDigest string `json:"pub_key_digest"`

	// UUID - device UUID (legacy identifier used by the remote access subsystem)
	UUID string `json:"uuid"`

	// System - system inventory.
	System SystemInfo `json:"system"`

	// LastReportedTime - timestamp of the last inventory report.
	LastReportedTime int64 `json:"last_reported_time"`

	// HeartbeatExpireSoft - timestamp when device will be considered delayed.
	HeartbeatExpireSoft int64 `json:"exp_soft"`

	// HeartbeatExpireHard - timestamp when device will be considered disconnected.
	HeartbeatExpireHard int64 `json:"exp_hard"`

	// Ancestors - array of ancestor nodes (incl. self).
	Ancestors []string `json:"ancestors"`

	// Attributes - device attributes.
	Attributes DeviceAttributes `json:"attributes"`

	// Title - device title.
	Title string `json:"title"`

	// AgentInterval - agent interval.
	// DEPRECATED: This field doesn't provide valid information anymore.
	AgentInterval string `json:"agentinterval"`

	// Status - device status (online, delayed or disconnected).
	Status string `json:"status"`

	// Tags - device tags.
	Tags []string `json:"tags"`
}

// InventoryListResponse contains list of device inventories.
type InventoryListResponse struct {
	Items []InventoryListItem `json:"items"`
	Total int                 `json:"total"`
}

// GetDeviceInventory returns the device inventory for the given device ID.
func (cli *Client) GetDeviceInventory(ctx context.Context, deviceID string) (*DeviceInventory, error) {
	deviceInventory := new(DeviceInventory)

	path := deviceInventoryPath + "/" + deviceID

	err := cli.Call(ctx, http.MethodGet, path, nil, deviceInventory)
	if err != nil {
		return nil, err
	}

	return deviceInventory, nil
}

// ListDeviceInventory returns a list of device inventories based on provided query.
func (cli *Client) ListDeviceInventory(ctx context.Context, query InventoryListQuery) (*InventoryListResponse, error) {
	queryParameters, err := query.String()
	if err != nil {
		return nil, fmt.Errorf("failed to encode query: %w", err)
	}

	requestPath := deviceInventoryPath + "?" + queryParameters

	response := new(InventoryListResponse)

	if err = cli.Call(ctx, http.MethodGet, requestPath, nil, response); err != nil {
		return nil, err
	}

	return response, nil
}
