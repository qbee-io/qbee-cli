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
	"context"
	"errors"
	"fmt"
	"net/http"
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

// Login returns a new authenticated API Client.
func (cli *Client) Login(ctx context.Context, email, password string) (string, error) {
	request := &LoginRequest{
		Email:    email,
		Password: password,
	}

	response := new(LoginResponse)

	if err := cli.Call(ctx, http.MethodPost, loginPath, request, &response); err != nil {
		if apiError := make(Error); errors.As(err, &apiError) {
			// Two-factor authentication is unsupported, so let's return a meaningful error message.
			if _, has2FAChallenge := apiError["challenge"].(string); has2FAChallenge {
				return "", fmt.Errorf("two-factor authentication is unsupported")
			}
		}

		return "", err
	}

	return response.Token, nil
}
