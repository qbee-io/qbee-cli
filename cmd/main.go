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
	"path/filepath"
	"time"

	"github.com/qbee-io/qbee-cli/client"
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
	},
}

// Command represents a command line command.

type Config struct {
	AuthToken string `json:"token"`
	BaseURL   string `json:"base_url"`
}

func WriteConfig(config Config) error {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	qbeeConfigDir := filepath.Join(dirname, ".qbee")

	if err := os.MkdirAll(qbeeConfigDir, 0700); err != nil {
		return err
	}

	configFile := filepath.Join(qbeeConfigDir, "qbee-cli.json")

	jsonConfig, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(configFile, jsonConfig, 0600)

	return err
}

func ReadConfig() (*Config, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	qbeeConfigDir := filepath.Join(dirname, ".qbee")

	if err := os.MkdirAll(qbeeConfigDir, 0700); err != nil {
		return nil, err
	}

	configFile := filepath.Join(qbeeConfigDir, "qbee-cli.json")

	jsonConfig, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	config := new(Config)
	err = json.Unmarshal(jsonConfig, config)

	token, err := client.DecodeAccessToken(config.AuthToken, client.StandardClaims{})

	if err != nil {
		return nil, err
	}

	if token.Claims.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("token expired")
	}

	Infof("Using cached token with expiry: %s\n", time.Unix(token.Claims.ExpiresAt, 0).Format(time.RFC3339))
	return config, err
}

func GetAuthenticatedClient(ctx context.Context) (*client.Client, error) {
	if os.Getenv("QBEE_EMAIL") != "" && os.Getenv("QBEE_PASSWORD") != "" {
		email := os.Getenv("QBEE_EMAIL")
		password := os.Getenv("QBEE_PASSWORD")
		cli := client.New()
		if os.Getenv("QBEE_BASEURL") != "" {
			cli = cli.WithBaseURL(os.Getenv("QBEE_BASEURL"))
		}
		if err := cli.Authenticate(ctx, email, password); err != nil {
			return nil, err
		}
		Warnf("Using credentials from environment variables QBEE_EMAIL and QBEE_PASSWORD. Consider doing a qbee-cli login.")
		return cli, nil
	}

	if config, err := ReadConfig(); err == nil {
		return client.New().WithBaseURL(config.BaseURL).WithAuthToken(config.AuthToken), nil
	}
	return nil, fmt.Errorf("no authentication mechanism found")
}
