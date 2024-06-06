// Copyright 2024 qbee.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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
	"go.qbee.io/client/config"
)

var configCommand = Command{
	Description: "Config commands",
	SubCommands: map[string]Command{
		"save":   configSaveCommand,
		"show":   configShowCommand,
		"commit": configCommitCommand,
	},
}

const (
	configSaveConfigOption     = "config"
	configNodeIDOption         = "node"
	configTagIDOption          = "tag"
	configSaveBundleNameOption = "bundle"
	configCommitMessageOption  = "commit-message"
	configSaveTemplateParams   = "template-parameters"
	configShowScopeOption      = "scope"
)

var configCommitCommand = Command{
	Description: "Commit a configuration",
	Options: []Option{
		{
			Name:     configCommitMessageOption,
			Short:    "m",
			Help:     "Commit message",
			Required: true,
		},
	},
	Target: func(opts Options) error {
		ctx := context.Background()

		cli, err := client.LoginGetAuthenticatedClient(ctx)

		if err != nil {
			return err
		}

		configManager := client.NewConfigurationManager().
			WithClient(cli)

		err = configManager.Commit(
			ctx,
			opts[configCommitMessageOption],
		)

		if err != nil {
			return err
		}

		fmt.Println("Configuration committed successfully")
		return nil
	},
}

var configSaveCommand = Command{
	Description: "Apply a configuration",
	Options: []Option{
		{
			Name:     configSaveConfigOption,
			Short:    "c",
			Help:     "Configuration file",
			Required: true,
		},
		{
			Name:  configNodeIDOption,
			Short: "n",
			Help:  "Node ID (group or device)",
		},
		{
			Name:  configTagIDOption,
			Short: "t",
			Help:  "Tag name",
		},
		{
			Name:     configSaveBundleNameOption,
			Short:    "b",
			Help:     "Bundle name",
			Required: true,
		},
		{
			Name:  configSaveTemplateParams,
			Short: "p",
			Help:  "Template parameters",
		},
	},
	OptionsHandler: func(opts Options) error {
		if opts[configNodeIDOption] == "" && opts[configTagIDOption] == "" {
			return fmt.Errorf("either node or tag must be specified")
		}

		if opts[configNodeIDOption] != "" && opts[configTagIDOption] != "" {
			return fmt.Errorf("node and tag cannot be specified at the same time")
		}

		return nil
	},
	Target: func(opts Options) error {
		ctx := context.Background()

		cli, err := client.LoginGetAuthenticatedClient(ctx)

		if err != nil {
			return err
		}

		configManager := client.NewConfigurationManager().
			WithClient(cli)

		if opts[configSaveTemplateParams] != "" {
			templateParams := make(map[string]string)

			params := strings.Split(opts[configSaveTemplateParams], ",")

			for _, param := range params {
				kv := strings.Split(param, "=")
				if len(kv) != 2 {
					return fmt.Errorf("invalid template parameter: %s", param)
				}
				templateParams[kv[0]] = kv[1]
			}

			configManager.WithTemplateParameters(templateParams)
		}

		target := opts[configNodeIDOption]
		configManager.WithEntityType(config.EntityTypeNode)

		if opts[configTagIDOption] != "" {
			target = opts[configTagIDOption]
			configManager.WithEntityType(config.EntityTypeTag)
		}

		err = configManager.Save(
			ctx,
			target,
			opts[configSaveBundleNameOption],
			opts[configSaveConfigOption],
		)

		if err != nil {
			return err
		}

		fmt.Println("Configuration applied successfully")
		return nil
	},
}

var configShowCommand = Command{
	Description: "Show a configuration",
	Options: []Option{
		{
			Name:  configNodeIDOption,
			Short: "n",
			Help:  "Node ID",
		},
		{
			Name:  configTagIDOption,
			Short: "t",
			Help:  "Tag name",
		},
		{
			Name:  configShowScopeOption,
			Help:  "Configuration scope",
			Short: "s",
		},
	},
	OptionsHandler: func(opts Options) error {
		if opts[configNodeIDOption] == "" && opts[configTagIDOption] == "" {
			return fmt.Errorf("either node or tag must be specified")
		}

		if opts[configNodeIDOption] != "" && opts[configTagIDOption] != "" {
			return fmt.Errorf("node and tag cannot be specified at the same time")
		}

		return nil
	},
	Target: func(opts Options) error {

		ctx := context.Background()

		cli, err := client.LoginGetAuthenticatedClient(ctx)

		if err != nil {
			return err
		}

		configManager := client.NewConfigurationManager().
			WithClient(cli)

		target := opts[configNodeIDOption]
		configManager.WithEntityType(config.EntityTypeNode)
		if opts[configTagIDOption] != "" {
			target = opts[configTagIDOption]
			configManager.WithEntityType(config.EntityTypeTag)
		}

		configManager.WithEntityConfigScope(config.EntityConfigScopeAll)
		if opts[configShowScopeOption] != "" {
			configManager.WithEntityConfigScope(config.EntityConfigScope(opts[configShowScopeOption]))
		}

		config, err := configManager.GetConfig(ctx, target)

		if err != nil {
			return err
		}

		return json.NewEncoder(os.Stdout).Encode(config)
	},
}
