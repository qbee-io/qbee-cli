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

// RemoteAccessTarget defines
type RemoteAccessTarget struct {
	// Protocol is the protocol used for the remote access.
	// Can be either "tcp" or "udp".
	Protocol string

	// RemoteHost is the host of the remote machine to which the local port is forwarded.
	RemoteHost string

	// LocalPort is the port on the local machine to which the remote port is forwarded.
	LocalPort string

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
// <local_port>:<remote_host>:<remote_port>[/udp]
func ParseRemoteAccessTarget(targetString string) (RemoteAccessTarget, error) {
	target := RemoteAccessTarget{}

	// Split the target string into its parts.
	parts := strings.Split(targetString, ":")

	if len(parts) != 3 {
		return target, fmt.Errorf("invalid format")
	}

	var err error

	if target.LocalPort, err = parseNetworkPort(parts[0]); err != nil {
		return target, fmt.Errorf("invalid local port: %w", err)
	}

	target.RemoteHost = parts[1]

	// Only localhost is supported as remote host at the moment.
	// Once we roll out the new remote access solution, will allow any host in the same network.
	if target.RemoteHost != "localhost" {
		return target, fmt.Errorf("invalid remote host: only localhost is supported")
	}

	remotePort := parts[2]
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
	base := fmt.Sprintf("%s:%s:%s", target.LocalPort, target.RemoteHost, target.RemotePort)

	if target.Protocol == UDP {
		return base + "/udp"
	}

	return base
}

// parseNetworkPort parses a network port string and accepts either a port number or "stdio".
func parseNetworkPort(portString string) (string, error) {
	if portString == "" {
		return "", fmt.Errorf("empty port")
	}

	if portString == "stdio" {
		return portString, nil
	}

	_, err := strconv.ParseUint(portString, 10, 16)
	if err != nil {
		return "", fmt.Errorf("invalid port number")
	}

	return portString, nil
}
