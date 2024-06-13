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

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"go.qbee.io/client"
)

const (
	connectDeviceOption     = "device"
	connectTargetOption     = "target"
	connectConfigFileOption = "config"
	connectAllowFailures    = "allow-failures"
	connectShellOption      = "shell"
	connectCommandOption    = "command"
)

// Example config file:
// [
//   {
//     "device_id": "2c40ccd4f8587e3e98d66e9092db4e8b0827906cb79c764b77cb2090b9acb7c8",
//     "targets": [
//       "5050:localhost:22",
//       "5151:localhost:80"
//     ]
//   },
//   {
//     "device_id": "09790f5388e3180793192fa6952e4c13c25bee650ddbf707d2764ac349d54046",
//     "targets": [
//       "8080:localhost:22",
//       "8081:localhost:80"
//     ]
//   }
// ]

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
		{
			Name: connectAllowFailures,
			Help: "Allow one or more failures",
			Flag: "true",
		},
		{
			Name: connectShellOption,
			Help: "Start a shell",
			Flag: "true",
		},
		{
			Name: connectCommandOption,
			Help: "Run a command",
		},
	},
	OptionsHandler: func(opts Options) error {
		if opts[connectConfigFileOption] != "" {
			return nil
		}

		if opts[connectDeviceOption] == "" {
			return fmt.Errorf("missing device ID")
		}

		if opts[connectShellOption] == "true" {
			return nil
		}

		if opts[connectTargetOption] == "" {
			return fmt.Errorf("missing target")
		}

		return nil
	},
	Target: func(opts Options) error {

		ctx := context.Background()

		cli, err := client.LoginGetAuthenticatedClient(ctx)
		if err != nil {
			return err
		}

		if opts[connectShellOption] == "true" {
			return cli.ConnectShell(ctx, opts[connectDeviceOption], opts[connectCommandOption])
		}

		if opts[connectConfigFileOption] != "" {
			remoteAccessTargets := make([]client.RemoteAccessConnection, 0)
			configBytes, err := os.ReadFile(opts[connectConfigFileOption])
			if err != nil {
				return fmt.Errorf("error reading config file: %w", err)
			}

			if err := json.Unmarshal(configBytes, &remoteAccessTargets); err != nil {
				return fmt.Errorf("error parsing config file: %w", err)
			}

			if len(remoteAccessTargets) == 0 {
				return fmt.Errorf("no connections defined in config file")
			}
			return cli.ConnectMulti(ctx, remoteAccessTargets, opts[connectAllowFailures] == "true")
		}

		return cli.ParseConnect(ctx,
			opts[connectDeviceOption],
			strings.Split(opts[connectTargetOption], ","))
	},
}
