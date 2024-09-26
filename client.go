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
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// DefaultBaseURL is the default base-URL used by the client.
const DefaultBaseURL = "https://www.app.qbee.io"

// Client encapsulates communication with the public API.
type Client struct {
	baseURL      string
	authToken    string
	httpClient   *http.Client
	refreshToken string
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

// GetHTTPClient returns the HTTP client used by the client.
func (cli *Client) GetHTTPClient() *http.Client {
	return cli.httpClient
}

// WithBaseURL sets the base URL of the API endpoint.
func (cli *Client) WithBaseURL(baseURL string) *Client {
	cli.baseURL = strings.TrimSuffix(baseURL, "/")
	return cli
}

// GetBaseURL returns the base URL used by the client.
func (cli *Client) GetBaseURL() string {
	return cli.baseURL
}

// WithAuthToken sets the authentication token for the client.
func (cli *Client) WithAuthToken(authToken string) *Client {
	cli.authToken = authToken
	return cli
}

// WithRefreshToken sets the refresh token for the client.
func (cli *Client) WithRefreshToken(refreshToken string) *Client {
	cli.refreshToken = refreshToken
	return cli
}

// GetAuthToken returns the authentication token used by the client.
func (cli *Client) GetAuthToken() string {
	return cli.authToken
}

// GetRefreshToken returns the refresh token used by the client.
func (cli *Client) GetRefreshToken() string {
	return cli.refreshToken
}

// SetRefreshToken sets the refresh token used by the client.
func (cli *Client) SetRefreshToken(refreshToken string) {
	cli.refreshToken = refreshToken
}

// Request sends an HTTP request to the API and returns the HTTP response.
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
	if response, err = cli.DoWithRefresh(request); err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	return response, nil
}

func (cli *Client) DoWithRefresh(request *http.Request) (*http.Response, error) {

	var bodyBuffer bytes.Buffer
	if request.Body != nil {
		request.Body = io.NopCloser(io.TeeReader(request.Body, &bodyBuffer))
	}

	response, err := cli.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusUnauthorized {
		return response, nil
	}

	io.Copy(io.Discard, response.Body)
	response.Body.Close()

	request.Body = io.NopCloser(&bodyBuffer)

	if err := cli.RefreshToken(request.Context()); err != nil {
		return nil, fmt.Errorf("error refreshing token: %w", err)
	}

	request.Header.Set("Authorization", "Bearer "+cli.authToken)

	response, err = cli.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= http.StatusBadRequest {
		defer response.Body.Close()
		var responseBody []byte

		if responseBody, err = io.ReadAll(response.Body); err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}

		if len(responseBody) > 0 {
			return nil, ParseErrorResponse(responseBody)
		}

		return nil, fmt.Errorf("got an http error with no body: %d", response.StatusCode)
	}

	return response, nil
}

// Call the API using provided method and path.
// If src is not nil, it will be encoded as JSON and sent as the request body.
// If dst is not nil, the response body will be decoded as JSON into dst
func (cli *Client) Call(ctx context.Context, method, path string, src, dst any) error {
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

	for _, cookie := range response.Cookies() {
		if cookie.Name == refreshTokenCookieName {
			cli.SetRefreshToken(cookie.Value)
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

const refreshAuthTokenPath = "/api/v2/refresh-jwt"
const refreshTokenCookieName = "PHPSESSID"

// RefreshToken refreshes the client's authentication token.
func (cli *Client) RefreshToken(ctx context.Context) error {
	if cli.refreshToken == "" {
		return errors.New("no refresh token set")
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, cli.GetBaseURL()+refreshAuthTokenPath, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	httpRequest.AddCookie(&http.Cookie{
		Name:  refreshTokenCookieName,
		Value: cli.refreshToken,
	})

	var httpResponse *http.Response
	if httpResponse, err = cli.GetHTTPClient().Do(httpRequest); err != nil {
		return fmt.Errorf("error refreshing token: %w", err)
	}

	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		var responseBody []byte
		responseBody, _ = io.ReadAll(httpResponse.Body)
		return ParseErrorResponse(responseBody)
	}

	cli.WithAuthToken(httpResponse.Header.Get("Refreshed-Token"))

	newLoginConfig := LoginConfig{
		BaseURL:      cli.GetBaseURL(),
		AuthToken:    cli.GetAuthToken(),
		RefreshToken: cli.GetRefreshToken(),
	}

	if err := LoginWriteConfig(newLoginConfig); err != nil {
		return fmt.Errorf("error saving refresh token: %w", err)
	}
	return nil
}
