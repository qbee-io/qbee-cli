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

package client

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const loginPath = "/api/v2/login"

// LoginRequest is the request body for the Login API.
type LoginRequest struct {
	// Email is the user email.
	Email string `json:"email"`

	// Password is the user password.
	Password string `json:"password"`
}

// LoginResponse is the response body for the Login API.
type LoginResponse struct {
	// Token is the authentication token to be used as Bearer token in the Authorization header.
	Token string `json:"token"`
}

// LoginConfig is the configuration file for the CLI authentication.
type LoginConfig struct {
	AuthToken    string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	BaseURL      string `json:"base_url"`
}

// LoginWriteConfig writes the CLI authentication configuration to the user's home directory.
func LoginWriteConfig(config LoginConfig) error {
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

// LoginReadConfig reads the CLI authentication configuration from the user's home directory.
func LoginReadConfig() (*LoginConfig, error) {
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

	config := new(LoginConfig)
	if err := json.Unmarshal(jsonConfig, config); err != nil {
		return nil, err
	}

	return config, nil
}

// LoginGetAuthenticatedClient returns a new authenticated API Client.
func LoginGetAuthenticatedClient(ctx context.Context) (*Client, error) {
	if os.Getenv("QBEE_EMAIL") != "" && os.Getenv("QBEE_PASSWORD") != "" {
		email := os.Getenv("QBEE_EMAIL")
		password := os.Getenv("QBEE_PASSWORD")
		cli := New()
		if os.Getenv("QBEE_BASEURL") != "" {
			cli = cli.WithBaseURL(os.Getenv("QBEE_BASEURL"))
		}
		if err := cli.Authenticate(ctx, email, password); err != nil {
			return nil, err
		}

		return cli, nil
	}

	var config *LoginConfig
	var err error
	if config, err = LoginReadConfig(); err != nil {
		return nil, err
	}

	client := New().WithBaseURL(config.BaseURL).
		WithAuthToken(config.AuthToken).
		WithRefreshToken(config.RefreshToken)

	return client, nil
}

// Login returns a new authenticated API Client.
func (cli *Client) Login(ctx context.Context, email, password string) (string, error) {
	request := &LoginRequest{
		Email:    email,
		Password: password,
	}

	response := new(LoginResponse)

	if err := cli.Call(ctx, http.MethodPost, loginPath, request, &response); err != nil {

		// If the error is an API error, check if it's a 2FA challenge.
		if apiError := make(Error); errors.As(err, &apiError) {

			if challenge, has2FAChallenge := apiError["challenge"].(string); has2FAChallenge {
				return cli.Login2FA(ctx, challenge)
			}
			return "", err
		}
		return "", err
	}

	return response.Token, nil
}

// Login2FAMethod is a container for the 2FA provider and code used during login.
type Login2FAMethod struct {
	Provider string
	Code     string
}

// Login2FARequest contains the request body for the Login 2FA API.
type Login2FARequest struct {
	Challenge string `json:"challenge,omitempty"`
	Provider  string `json:"preferProvider,omitempty"`
	Code      string `json:"code,omitempty"`
}

// Login2FAResponse is the response body for the Login 2FA API.
type Login2FAResponse struct {
	Challenge string `json:"challenge,omitempty"`
	Token     string `json:"token,omitempty"`
}

var valid2FAProviders = []string{"google", "email"}

const login2FAChallengeGetPath = "/api/v2/challenge-get"
const login2FAChallengeVerifyPath = "/api/v2/challenge-verify"

// Login2FA returns a new authenticated API Client.
func (cli *Client) Login2FA(ctx context.Context, challenge string) (string, error) {
	login2FAMethod := findProvided2FAMethod()
	if login2FAMethod == nil {
		var err error
		provider, err := prompt2FAProvider()
		if err != nil {
			return "", err
		}

		login2FAMethod = &Login2FAMethod{
			Provider: provider,
		}
	}

	challengeGetRequest := &Login2FARequest{
		Challenge: challenge,
		Provider:  login2FAMethod.Provider,
	}
	challengeGetResponse := new(Login2FAResponse)
	if err := cli.Call(ctx, http.MethodPost, login2FAChallengeGetPath, challengeGetRequest, &challengeGetResponse); err != nil {
		return "", err
	}

	// The code might already have been provided as an environment variable.
	if login2FAMethod.Code == "" {
		code, err := prompt2FACode()
		if err != nil {
			return "", err
		}

		login2FAMethod.Code = code
	}

	challengeVerifyRequest := &Login2FARequest{
		Challenge: challengeGetResponse.Challenge,
		Code:      login2FAMethod.Code,
	}
	challengeVerifyResponse := new(Login2FAResponse)
	if err := cli.Call(ctx, http.MethodPost, login2FAChallengeVerifyPath, challengeVerifyRequest, &challengeVerifyResponse); err != nil {
		return "", err
	}

	return challengeVerifyResponse.Token, nil
}

func findProvided2FAMethod() *Login2FAMethod {
	if os.Getenv("QBEE_2FA_CODE") != "" {
		fmt.Printf("Using 2FA code from environment variable QBEE_2FA_CODE as a google 2FA provider\n")

		return &Login2FAMethod{
			Provider: "google",
			Code:     os.Getenv("QBEE_2FA_CODE"),
		}
	}

	return nil
}

func prompt2FAProvider() (string, error) {
	fmt.Printf("Select 2FA provider:\n")

	for i, provider := range valid2FAProviders {
		fmt.Printf("%d) %s\n", i+1, provider)
	}

	fmt.Printf("Choice: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}

	providerIndex := scanner.Text()

	index, err := strconv.Atoi(providerIndex)
	if err != nil {
		return "", err
	}

	if index < 1 || index > len(valid2FAProviders) {
		return "", fmt.Errorf("invalid provider")
	}

	provider := valid2FAProviders[index-1]
	return provider, nil
}

func prompt2FACode() (string, error) {
	fmt.Printf("Enter 2FA code: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		return "", err
	}

	code := scanner.Text()
	return code, nil
}
