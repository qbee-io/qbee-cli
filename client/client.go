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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// DefaultBaseURL is the default base-URL used by the client.
const DefaultBaseURL = "https://www.app.qbee.io"

// Client encapsulates communication with the public API.
type Client struct {
	baseURL    string
	authToken  string
	httpClient *http.Client
}

// New returns a new instance of a public API client.
func New() *Client {
	return &Client{
		baseURL:    DefaultBaseURL,
		httpClient: http.DefaultClient,
	}
}

// WithHTTPClient sets the HTTP client to use for requests.
func (cli *Client) WithHTTPClient(httpClient *http.Client) *Client {
	cli.httpClient = httpClient
	return cli
}

// WithBaseURL sets the base URL of the API endpoint.
func (cli *Client) WithBaseURL(baseURL string) *Client {
	cli.baseURL = baseURL
	return cli
}

// WithAuthToken sets the authentication token for the client.
func (cli *Client) WithAuthToken(authToken string) *Client {
	cli.authToken = authToken
	return cli
}

// Request performs an HTTP request to the API.
func (cli *Client) Request(ctx context.Context, method, path string, src any) (*http.Response, error) {
	if !strings.HasPrefix(path, "/") {
		return nil, fmt.Errorf("path %s must start with /", path)
	}

	var body io.ReadWriter

	if src != nil {
		body = new(bytes.Buffer)

		if err := json.NewEncoder(body).Encode(src); err != nil {
			return nil, fmt.Errorf("error encoding JSON request: %w", err)
		}
	}

	url := cli.baseURL + path

	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	if cli.authToken != "" {
		request.Header.Set("Authorization", "Bearer "+cli.authToken)
	}

	request.Header.Set("User-Agent", UserAgent)

	var response *http.Response
	if response, err = cli.httpClient.Do(request); err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	return response, nil
}

// request sends an HTTP request with optional JSON payload (src) and optionally decodes JSON response to dst.
func (cli *Client) request(ctx context.Context, method, path string, src, dst any) error {
	response, err := cli.Request(ctx, method, path, src)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var responseBody []byte
	if responseBody, err = io.ReadAll(response.Body); err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if response.StatusCode >= http.StatusBadRequest {
		if len(responseBody) > 0 {
			return ParseErrorResponse(responseBody)
		}

		return fmt.Errorf("got an http error with no body: %d", response.StatusCode)
	}

	if dst != nil && len(responseBody) > 0 {
		if err = json.Unmarshal(responseBody, dst); err != nil {
			return fmt.Errorf("error decoding JSON response (%w): %s", err, responseBody)
		}
	}

	return nil
}

// Authenticate authenticates the client instance with the given email and password.
func (cli *Client) Authenticate(ctx context.Context, email string, password string) error {
	token, err := cli.Login(ctx, email, password)
	if err != nil {
		return fmt.Errorf("error authenticating: %w", err)
	}

	cli.WithAuthToken(token)

	return nil
}
