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

const (
	TCP = "tcp"
	UDP = "udp"
)

// RemoteAccessTarget defines
type RemoteAccessTarget struct {
	// Protocol is the protocol used for the remote access.
	// Can be either "tcp" or "udp".
	Protocol string

	// Device identifies the target device by its public key digest (SHA256).
	Device string

	// LocalPort is the port on the local machine to which the remote port is forwarded.
	// If set to 0, a random port will be chosen.
	LocalPort uint16

	// RemotePort is the port on the remote machine to which the local port is forwarded.
	RemotePort uint16
}

// isValidDeviceID checks if the provided device ID is valid.
func isValidDeviceID(deviceID string) bool {
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
// <proto>:<local_port>:<device>:<remote_port>
func ParseRemoteAccessTarget(targetString string) (RemoteAccessTarget, error) {
	target := RemoteAccessTarget{}

	// Split the target string into its parts.
	parts := strings.Split(targetString, ":")

	if len(parts) != 4 {
		return target, fmt.Errorf("invalid format")
	}

	switch parts[0] {
	case TCP, UDP:
		target.Protocol = parts[0]
	default:
		return target, fmt.Errorf("invalid protocol")
	}

	var err error

	if target.LocalPort, err = parseNetworkPort(parts[1]); err != nil {
		return target, fmt.Errorf("invalid local port: %w", err)
	}

	if target.Device = parts[2]; !isValidDeviceID(target.Device) {
		return target, fmt.Errorf("invalid device")
	}

	if target.RemotePort, err = parseNetworkPort(parts[3]); err != nil {
		return target, fmt.Errorf("invalid remote port: %w", err)
	}

	return target, nil
}

// parseNetworkPort parses a network port string.
func parseNetworkPort(portString string) (uint16, error) {
	if portString == "" {
		return 0, fmt.Errorf("empty port")
	}

	portUint, err := strconv.ParseUint(portString, 10, 16)
	if err != nil {
		return 0, fmt.Errorf("invalid port number")
	}

	return uint16(portUint), nil
}
