// Copyright 2024 qbee.io
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
	"fmt"
	"runtime"

	"go.qbee.io/client"
)

const (
	shellDeviceOption  = "device"
	shellCommandOption = "command"
)

var shellCommand = Command{
	Description: "Start a shell on a device",
	Options: []Option{
		{
			Name:     shellDeviceOption,
			Short:    "d",
			Help:     "Device ID",
			Required: true,
		},
		{
			Name:     shellCommandOption,
			Short:    "c",
			Help:     "Command to execute as comma-separated list of arguments",
			Required: false,
		},
	},
	Target: func(opts Options) error {
		ctx := context.Background()
		cli, err := client.LoginGetAuthenticatedClient(ctx)
		if err != nil {
			return err
		}

		if runtime.GOOS == "windows" {
			return fmt.Errorf("shell is not supported on Windows")
		}

		deviceID := opts[shellDeviceOption]
		if err := cli.ConnectShell(ctx, deviceID, opts[shellCommandOption]); err != nil {
			return err
		}

		return nil
	},
}