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
	"time"
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
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
}

type LoginConfig struct {
	AuthToken string `json:"token"`
	BaseURL   string `json:"base_url"`
}

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
	err = json.Unmarshal(jsonConfig, config)

	token, err := DecodeAccessToken(config.AuthToken, StandardClaims{})

	if err != nil {
		return nil, err
	}

	if token.Claims.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("token expired")
	}

	//Infof("Using cached token with expiry: %s\n", time.Unix(token.Claims.ExpiresAt, 0).Format(time.RFC3339))
	return config, err
}

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

	if config, err := LoginReadConfig(); err == nil {
		return New().WithBaseURL(config.BaseURL).WithAuthToken(config.AuthToken), nil
	}
	return nil, fmt.Errorf("no authentication mechanism found")
}

// Login returns a new authenticated API Client.
func (cli *Client) Login(ctx context.Context, email, password string) (string, error) {
	request := &LoginRequest{
		Email:    email,
		Password: password,
	}

	response := new(LoginResponse)

	if err := cli.Call(ctx, http.MethodPost, loginPath, request, &response); err != nil {
		if apiError := make(Error); errors.As(err, &apiError) {
			if challenge, has2FAChallenge := apiError["challenge"].(string); has2FAChallenge {
				return cli.Login2FA(ctx, challenge)
			}
			return "", err
		}
	}

	return response.Token, nil
}

type Login2FARequest struct {
	Challenge string `json:"challenge"`
	Provider  string `json:"preferProvider"`
}

type Login2FAVerifyRequest struct {
	Challenge string `json:"challenge"`
	Code      string `json:"code"`
}

const login2FAChallengeGetPath = "/api/v2/challenge-get"
const login2FAChallengeVerifyPath = "/api/v2/challenge-verify"

// Login2FA returns a new authenticated API Client.
func (cli *Client) Login2FA(ctx context.Context, challenge string) (string, error) {
	fmt.Printf("Enter provider (1 Google authenticator, 2 Email): ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}

	provider := "google"
	requestPrepare := &Login2FARequest{
		Challenge: challenge,
		Provider:  provider,
	}

	responsePrepare := new(LoginResponse)
	if err := cli.Call(ctx, http.MethodPost, login2FAChallengeGetPath, requestPrepare, &responsePrepare); err != nil {
		if apiError := make(Error); errors.As(err, &apiError) {
			var has2FAChallenge bool
			if challenge, has2FAChallenge = apiError["challenge"].(string); !has2FAChallenge {
				return "", err
			}
		}
	}

	fmt.Printf("Enter code: ")
	scanner = bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err = scanner.Err()
	if err != nil {
		log.Fatal(err)
	}

	code := scanner.Text()

	requestVerify := &Login2FAVerifyRequest{
		Challenge: responsePrepare.Challenge,
		Code:      code,
	}

	response := new(LoginResponse)

	if err := cli.Call(ctx, http.MethodPost, login2FAChallengeVerifyPath, requestVerify, &response); err != nil {
		return "", err
	}

	return response.Token, nil

}
