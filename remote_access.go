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
	"fmt"
	"strconv"
	"strings"
)

// Network protocols supported by the remote access.
const (
	TCP = "tcp"
	UDP = "udp"
)

// RemoteAccessConnection defines a remote access connection.
type RemoteAccessConnection struct {
	// DeviceID is the device ID of the device to which the remote access connection belongs.
	DeviceID string `json:"device_id"`

	// Targets is a list of remote access targets.
	Targets []string `json:"targets"`
}

// RemoteAccessTarget defines
type RemoteAccessTarget struct {
	// Protocol is the protocol used for the remote access.
	// Can be either "tcp" or "udp".
	Protocol string

	// LocalHost is the address on which the local port is bound.
	// Set to "localhost" to bind to the loopback interface.
	LocalHost string

	// LocalPort is the port on the local machine to which the remote port is forwarded.
	LocalPort string

	// RemoteHost is the host of the remote machine to which the local port is forwarded.
	RemoteHost string

	// RemotePort is the port on the remote machine to which the local port is forwarded.
	RemotePort string
}

// IsValidDeviceID checks if the provided device ID is valid.
func IsValidDeviceID(deviceID string) bool {
	// make sure length is correct before we check the content
	if len(deviceID) != 64 {
		return false
	}

	// iterate over every character and make sure it is a valid hex character
	for _, c := range deviceID {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}

	return true
}

// ParseRemoteAccessTarget parses a remote access target string.
// The target string has the following format:
// [<local_host>:]<local_port>:<remote_host>:<remote_port>[/udp]
func ParseRemoteAccessTarget(targetString string) (RemoteAccessTarget, error) {
	target := RemoteAccessTarget{}

	// Split the target string into its parts.
	parts := strings.Split(targetString, ":")

	if len(parts) != 3 && len(parts) != 4 {
		return target, fmt.Errorf("invalid format")
	}

	var err error

	var localHost string
	var localPort string
	var remoteHost string
	var remotePort string

	if len(parts) == 3 {
		localHost = "localhost"
		localPort = parts[0]
		remoteHost = parts[1]
		remotePort = parts[2]
	} else {
		localHost = parts[0]
		localPort = parts[1]
		remoteHost = parts[2]
		remotePort = parts[3]
	}

	target.LocalHost = localHost
	target.RemoteHost = remoteHost

	if target.LocalPort, err = parseLocalNetworkPort(localPort); err != nil {
		return target, fmt.Errorf("invalid local port: %w", err)
	}

	if strings.HasSuffix(remotePort, "/udp") {
		target.Protocol = UDP
		remotePort = strings.TrimSuffix(remotePort, "/udp")
	} else {
		target.Protocol = TCP
	}

	if target.RemotePort, err = parseNetworkPort(remotePort); err != nil {
		return target, fmt.Errorf("invalid remote port: %w", err)
	}

	return target, nil
}

// String returns the string representation of a remote access target.
func (target RemoteAccessTarget) String() string {
	base := fmt.Sprintf("%s:%s:%s:%s", target.LocalHost, target.LocalPort, target.RemoteHost, target.RemotePort)

	if target.Protocol == UDP {
		return base + "/udp"
	}

	return base
}

// parseLocalNetworkPort parses a network port string and accepts either a port number or "stdio".
func parseLocalNetworkPort(portString string) (string, error) {
	if portString == "stdio" {
		return portString, nil
	}

	if _, err := parseNetworkPort(portString); err != nil {
		return "", err
	}

	return portString, nil
}

// parseNetworkPort parses a network port string and accepts a port number.
func parseNetworkPort(portString string) (string, error) {
	if portString == "" {
		return "", fmt.Errorf("empty port")
	}

	_, err := strconv.ParseUint(portString, 10, 16)
	if err != nil {
		return "", fmt.Errorf("invalid port number")
	}

	return portString, nil
}
