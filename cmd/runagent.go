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
	"fmt"

	"go.qbee.io/client"
)

const (
	runAgentDeviceIDOption      = "device-id"
	runAgentGroupIDOption       = "group-id"
	runAgentTagOption           = "tag"
	runAgentAllowFailuresOption = "allow-failures"
)

var runAgentCommand = Command{
	Description: "Run the agent",
	Options: []Option{
		{
			Name:  runAgentDeviceIDOption,
			Short: "d",
			Help:  "Device ID (as Public Key Digest)",
		},
		{
			Name:  runAgentGroupIDOption,
			Short: "g",
			Help:  "Group ID",
		},
		{
			Name:  runAgentTagOption,
			Short: "t",
			Help:  "Tag",
		},
		{
			Name: runAgentAllowFailuresOption,
			Help: "Allow failures",
			Flag: "true",
		},
	},
	Target: func(opts Options) error {

		ctx := context.Background()

		cli, err := client.LoginGetAuthenticatedClient(ctx)
		if err != nil {
			return err
		}

		runAgentManager := client.NewRunAgentManager().
			WithClient(cli).
			WithAllowFailures(opts[runAgentAllowFailuresOption] == "true")

		if opts[runAgentDeviceIDOption] != "" {
			return runAgentManager.RunAgentDevice(ctx, opts[runAgentDeviceIDOption])
		}

		if opts[runAgentGroupIDOption] != "" {
			fmt.Printf("Running agent for group %s\n", opts[runAgentGroupIDOption])
			return runAgentManager.RunAgentGroup(ctx, opts[runAgentGroupIDOption])
		}

		if opts[runAgentTagOption] != "" {
			fmt.Printf("Running agent for tag %s\n", opts[runAgentTagOption])
			return runAgentManager.RunAgentTag(ctx, opts[runAgentTagOption])
		}

		return fmt.Errorf("missing device ID, group ID or tag")
	},
}
