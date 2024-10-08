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
	"fmt"
	"os"
)

const (
	mainLogLevel = "log-level"
)

// Main is the main command of the agent.
var Main = Command{
	Description: "Qbee Agent Command-Line Tool",
	Options: []Option{
		{
			Name:    mainLogLevel,
			Short:   "l",
			Help:    "Logging level: DEBUG, INFO, WARNING or ERROR.",
			Default: "INFO",
		},
	},
	OptionsHandler: func(opts Options) error {
		switch opts[mainLogLevel] {
		case "DEBUG":
			SetLogLevel(DEBUG)
		case "INFO":
			SetLogLevel(INFO)
		case "WARNING":
			SetLogLevel(WARNING)
		case "ERROR":
			SetLogLevel(ERROR)
		}

		return nil
	},
	SubCommands: map[string]Command{
		"connect": connectCommand,
		"version": versionCommand,
		"login":   loginCommand,
		"files":   filemanagerCommand,
		"devices": devicesCommand,
		"config":  configCommand,
		"run":     runAgentCommand,
		"term":    terminalCommand,
	},
}

func main() {
	if err := Main.Execute(os.Args[1:], nil); err != nil {
		fmt.Fprintf(os.Stderr, "qbee-cli error: %s\n", err)
		os.Exit(1)
	}
}
