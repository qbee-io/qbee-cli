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

package client

import (
	"context"
	"net/http"
)

const secretPath = "/api/v2/secret"

// CreateSecretRequest is the request body for creating a secret.
type CreateSecretRequest struct {
	// Value is the plain text value of the secret.
	Value string `json:"value"`
}

// CreateSecretResponse is the response from the server when creating a secret.
type CreateSecretResponse struct {
	// SecretRef is the reference to the created secret.
	SecretRef string `json:"secret_ref"`
}

// CreateSecret stores the provided plain text value as a secret and returns its reference.
func (cli *Client) CreateSecret(ctx context.Context, value string) (string, error) {
	request := CreateSecretRequest{Value: value}
	response := new(CreateSecretResponse)

	if err := cli.Call(ctx, http.MethodPost, secretPath, request, response); err != nil {
		return "", err
	}

	return response.SecretRef, nil
}
