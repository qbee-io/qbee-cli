// Copyright 2026 qbee.io
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
	"os"
	"strings"

	"go.qbee.io/client"
)

var copyCommand = Command{
	Description: "Copy files to or from a device",
	Usage:       "copy [<device ID>:]<sourcePath> [<device ID>:]<destinationPath>",
	OptionsHandler: func(opts Options) error {
		if len(os.Args) != 4 {
			return fmt.Errorf("invalid arguments")
		}

		source, destination := os.Args[2], os.Args[3]

		sourceDevice, sourcePath, isRemote := strings.Cut(source, ":")
		if isRemote {
			if !client.IsValidDeviceID(sourceDevice) {
				return fmt.Errorf("invalid source device ID %s", sourceDevice)
			}
		} else {
			sourcePath = source
			sourceDevice = ""
		}

		destinationDevice, destinationPath, isRemote := strings.Cut(destination, ":")
		if isRemote {
			if sourceDevice != "" {
				return fmt.Errorf("both source and destination cannot be remote")
			}

			if !client.IsValidDeviceID(destinationDevice) {
				return fmt.Errorf("invalid destination device ID %s", destinationDevice)
			}
		} else {
			destinationPath = destination
			destinationDevice = ""
		}

		if sourceDevice == "" && destinationDevice == "" {
			return fmt.Errorf("either source or destination must be remote")
		}

		opts["sourceDevice"] = sourceDevice
		opts["sourcePath"] = sourcePath
		opts["destinationDevice"] = destinationDevice
		opts["destinationPath"] = destinationPath

		return nil
	},
	Target: func(opts Options) error {
		ctx := context.Background()
		cli, err := client.LoginGetAuthenticatedClient(ctx)
		if err != nil {
			return err
		}

		if opts["sourceDevice"] != "" {
			return cli.DownloadFileFromDevice(ctx, opts["sourceDevice"], opts["sourcePath"], opts["destinationPath"])
		}

		return cli.UploadFileToDevice(ctx, opts["destinationDevice"], opts["sourcePath"], opts["destinationPath"])
	},
}
