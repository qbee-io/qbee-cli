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
	"os"

	"github.com/qbee-io/qbee-cli/client"
)

const (
	connectTargetOption = "target"
)

var connectCommand = Command{
	Description: "Connect to a device",
	Options: []Option{
		{
			Name:     connectTargetOption,
			Short:    "t",
			Help:     "Target defined as <proto>:<localPort>:<deviceID>:<remotePort>",
			Required: true,
		},
	},
	Target: func(opts Options) error {
		email := os.Getenv("QBEE_EMAIL")
		password := os.Getenv("QBEE_PASSWORD")

		ctx := context.Background()

		target, err := client.ParseRemoteAccessTarget(opts[connectTargetOption])
		if err != nil {
			return fmt.Errorf("error parsing target: %w", err)
		}

		cli := client.New()
		if baseURL, ok := os.LookupEnv("QBEE_BASEURL"); ok {
			cli = cli.WithBaseURL(baseURL)
		}

		if err = cli.Authenticate(ctx, email, password); err != nil {
			return err
		}

		return cli.Connect(ctx, target)
	},
}
