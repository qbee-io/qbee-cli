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
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/qbee-io/qbee-cli/client"
)

const (
	connectDeviceOption     = "device"
	connectTargetOption     = "target"
	connectConfigFileOption = "config"
)

var RemoteAccessTargets = make([]client.RemoteAccessConnection, 0)

var connectCommand = Command{
	Description: "Connect to a device",
	Options: []Option{
		{
			Name:  connectDeviceOption,
			Short: "d",
			Help:  "Device ID (as Public Key Digest)",
		},
		{
			Name:  connectTargetOption,
			Short: "t",
			Help:  "Comma-separated targets definition <localPort>:<remoteHost>:<remotePort>[/udp]",
		},
		{
			Name:  connectConfigFileOption,
			Short: "c",
			Help:  "Config file to use",
		},
	},
	OptionsHandler: func(opts Options) error {
		if opts[connectConfigFileOption] != "" {
			configBytes, err := os.ReadFile(opts[connectConfigFileOption])
			if err != nil {
				return fmt.Errorf("error reading config file: %w", err)
			}

			err = json.Unmarshal(configBytes, &RemoteAccessTargets)
			if err != nil {
				return fmt.Errorf("error parsing config file: %w", err)
			}

			if len(RemoteAccessTargets) == 0 {
				return fmt.Errorf("no connections defined in config file")
			}
			return nil
		}
		if opts[connectDeviceOption] == "" {
			return fmt.Errorf("missing device ID")
		}
		if opts[connectTargetOption] == "" {
			return fmt.Errorf("missing target")
		}
		RemoteAccessTargets = append(RemoteAccessTargets, client.RemoteAccessConnection{
			DeviceID: opts[connectDeviceOption],
			Targets:  strings.Split(opts[connectTargetOption], ","),
		})
		return nil
	},
	Target: func(opts Options) error {
		email := os.Getenv("QBEE_EMAIL")
		password := os.Getenv("QBEE_PASSWORD")

		ctx := context.Background()

		remoteAccessMap := make(map[string][]client.RemoteAccessTarget, 0)

		for _, target := range RemoteAccessTargets {
			if !client.IsValidDeviceID(target.DeviceID) {
				return fmt.Errorf("invalid device ID %s", target.DeviceID)
			}

			targets := make([]client.RemoteAccessTarget, 0)
			for _, targetString := range target.Targets {
				target, err := client.ParseRemoteAccessTarget(targetString)
				if err != nil {
					return fmt.Errorf("error parsing target %s: %w", targetString, err)
				}

				targets = append(targets, target)
			}

			remoteAccessMap[target.DeviceID] = targets
		}

		cli := client.New()
		if baseURL, ok := os.LookupEnv("QBEE_BASEURL"); ok {
			cli = cli.WithBaseURL(baseURL)
		}

		if err := cli.Authenticate(ctx, email, password); err != nil {
			return err
		}
		return cli.Connect(ctx, remoteAccessMap)
	},
}
