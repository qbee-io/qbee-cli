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

package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/qbee-io/qbee-cli/client"
)

const (
	connectDeviceOption = "device"
	connectTargetOption = "target"
)

var connectCommand = Command{
	Description: "Connect to a device",
	Options: []Option{
		{
			Name:     connectDeviceOption,
			Short:    "d",
			Help:     "Device ID (as Public Key Digest)",
			Required: true,
		},
		{
			Name:     connectTargetOption,
			Short:    "t",
			Help:     "Comma-separated targets definition <localPort>:<remoteHost>:<remotePort>[/udp]",
			Required: true,
		},
	},
	Target: func(opts Options) error {
		ctx := context.Background()

		deviceID := opts[connectDeviceOption]
		if !client.IsValidDeviceID(deviceID) {
			return fmt.Errorf("invalid device ID %s", deviceID)
		}

		targets := make([]client.RemoteAccessTarget, 0)
		for _, targetString := range strings.Split(opts[connectTargetOption], ",") {
			target, err := client.ParseRemoteAccessTarget(targetString)
			if err != nil {
				return fmt.Errorf("error parsing target %s: %w", targetString, err)
			}

			targets = append(targets, target)
		}

		cli, err := GetAuthenticatedClient(ctx)
		if err != nil {
			return err
		}

		return cli.Connect(ctx, deviceID, targets)
	},
}
