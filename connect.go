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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	chisel "github.com/jpillora/chisel/client"
	"golang.org/x/net/http/httpproxy"
)

const (
	remoteAccessKeepAlive     = 25 * time.Second
	remoteAccessLegacyURL     = "%s/qbee-connect/%s/"
	remoteAccessInterfaceName = "qbee0"
)

// RemoteAccessConnectionInfo contains parameters used to establish a remote access connection.
type RemoteAccessConnectionInfo struct {
	// URL of the remote access server. If empty, a legacy URL will be used.
	URL string `json:"url"`

	// VPNIndex is the legacy VPN index to use for the connection.
	VPNIndex string `json:"vpn_idx"`

	// Username is the username to use for authentication.
	Username string `json:"username"`

	// Password is the password to use for authentication.
	Password string `json:"password"`

	// Address is the internal device address client is connecting to.
	Address string `json:"address"`
}

// RemoteAccessToken contains remote console
type RemoteAccessToken struct {
	// Title is the title of the device.
	Title string `json:"title"`

	// Connection contains the connection details.
	Connection RemoteAccessConnectionInfo `json:"connection"`
}

// AuthString returns the authentication string for the remote access token.
func (token RemoteAccessToken) AuthString() string {
	return fmt.Sprintf("%s:%s", token.Connection.Username, token.Connection.Password)
}

// ProxyURL returns the proxy URL for the remote access token.
func (token RemoteAccessToken) ProxyURL() (string, error) {
	serverURL, err := url.Parse(token.Connection.URL)
	if err != nil {
		return "", fmt.Errorf("error parsing remote access server URL: %w", err)
	}

	var proxyURL *url.URL
	if proxyURL, err = httpproxy.FromEnvironment().ProxyFunc()(serverURL); err != nil {
		return "", fmt.Errorf("error getting proxy URL: %w", err)
	}

	if proxyURL == nil {
		return "", nil
	}

	return proxyURL.String(), nil
}

// RemoteAccessTokenResponse is the response returned by the API when requesting a remote console token.
type RemoteAccessTokenResponse map[string]RemoteAccessToken

// RemoteAccessTokenRequest is the request sent to the API when requesting a remote console token.
type RemoteAccessTokenRequest struct {
	// DeviceID is the PublicKeyDigest of the device for which the token is requested.
	// Required.
	DeviceID string

	// Application is the name of the application for which the token is requested.
	// Optional - defaults to "qbee-cli".
	Application string

	// Username is the username for which the token is requested.
	// Setting this value will allow to identify the user in the audit log.
	// Optional.
	Username string

	// Ports is the list of ports for which the token is requested.
	// Setting this value will allow to identify requested ports in the audit log.
	// Optional.
	Ports []string
}

// remoteAccessTokenV2Path is the path of the remote access token API.
// First argument is the device ID, second argument is the URL-encoded query string.
const remoteAccessTokenV2Path = "/api/v2/remoteconsoletokenv2/%s?%s"

// GetRemoteAccessToken returns a remote console token for the specified device.
func (cli *Client) GetRemoteAccessToken(ctx context.Context, req RemoteAccessTokenRequest) (*RemoteAccessToken, error) {
	if req.DeviceID == "" {
		return nil, fmt.Errorf("device ID is required")
	}

	if req.Application == "" {
		req.Application = Name
	}

	urlValues := url.Values{}
	urlValues.Set("app_name", req.Application)

	if req.Username != "" {
		urlValues.Set("username", req.Username)
	}

	if len(req.Ports) > 0 {
		jsonEncodedPorts, err := json.Marshal(req.Ports)
		if err != nil {
			return nil, fmt.Errorf("error encoding ports: %w", err)
		}

		urlValues.Set("ports", string(jsonEncodedPorts))
	}

	path := fmt.Sprintf(remoteAccessTokenV2Path, req.DeviceID, urlValues.Encode())

	response := make(RemoteAccessTokenResponse)

	if err := cli.Call(ctx, http.MethodGet, path, nil, &response); err != nil {
		return nil, err
	}

	remoteAccessToken := response[req.DeviceID]

	// for legacy remote access servers, set the URL based on the VPN index
	if remoteAccessToken.Connection.URL == "" {
		remoteAccessToken.Connection.URL = fmt.Sprintf(
			remoteAccessLegacyURL, cli.baseURL, remoteAccessToken.Connection.VPNIndex)
	}

	return &remoteAccessToken, nil
}

// ConnectMulti establishes connections to multiple remote devices concurrently.
func (cli *Client) ConnectMulti(ctx context.Context, connections []RemoteAccessConnection, allowFailures bool) error {
	wg := sync.WaitGroup{}
	errChan := make(chan error)
	done := make(chan bool)

	for _, conn := range connections {
		wg.Add(1)

		go func(connection RemoteAccessConnection) {
			defer wg.Done()

			if err := cli.ParseConnect(ctx, connection.DeviceID, connection.Targets); err != nil {
				errChan <- fmt.Errorf("error connecting to device %s: %w", connection.DeviceID, err)
			}
		}(conn)
	}

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		done <- true
	}(&wg)

	for {
		select {
		case <-done:
			return nil
		case err := <-errChan:
			if !allowFailures {
				return err
			}
			fmt.Printf("%s\n", err)
		}
	}
}

// ParseConnect parses a device ID and a list of targets and establishes a connection to the device.
func (cli *Client) ParseConnect(ctx context.Context, deviceID string, targets []string) error {

	if !IsValidDeviceID(deviceID) {
		return fmt.Errorf("invalid device ID %s", deviceID)
	}

	parsedTargets := make([]RemoteAccessTarget, 0)

	for _, targetString := range targets {
		target, err := ParseRemoteAccessTarget(targetString)
		if err != nil {
			return fmt.Errorf("error parsing target %s: %w", targetString, err)
		}

		parsedTargets = append(parsedTargets, target)
	}

	if len(parsedTargets) == 0 {
		return fmt.Errorf("no targets defined for device %s", deviceID)
	}

	return cli.Connect(ctx, deviceID, parsedTargets)
}

// Connect establishes a connection to a remote device.
func (cli *Client) Connect(ctx context.Context, deviceID string, targets []RemoteAccessTarget) error {
	ports := make([]string, len(targets))
	for _, target := range targets {
		ports = append(ports, fmt.Sprintf("%s:%s", target.Protocol, target.RemotePort))
	}

	remoteAccessTokenRequest := RemoteAccessTokenRequest{
		DeviceID: deviceID,
		Ports:    ports,
	}

	remoteAccessToken, err := cli.GetRemoteAccessToken(ctx, remoteAccessTokenRequest)
	if err != nil {
		return err
	}

	var proxyURL string
	if proxyURL, err = remoteAccessToken.ProxyURL(); err != nil {
		return fmt.Errorf("error getting proxy URL: %w", err)
	}

	address := remoteAccessToken.Connection.Address
	if address == "" {
		var deviceInventory *DeviceInventory
		if deviceInventory, err = cli.GetDeviceInventory(ctx, deviceID); err != nil {
			return fmt.Errorf("error getting device inventory: %w", err)
		}

		address = deviceInventory.SystemInfo.IPv4[remoteAccessInterfaceName]
	}

	remotes := make([]string, len(targets))
	for i, target := range targets {
		target.RemoteHost = address
		remotes[i] = target.String()
	}

	chiselClientConfig := &chisel.Config{
		Auth:      remoteAccessToken.AuthString(),
		Server:    remoteAccessToken.Connection.URL,
		Proxy:     proxyURL,
		KeepAlive: remoteAccessKeepAlive,
		Remotes:   remotes,
	}

	var chiselClient *chisel.Client
	if chiselClient, err = chisel.NewClient(chiselClientConfig); err != nil {
		return fmt.Errorf("error initializing remote access client: %w", err)
	}

	chiselClient.Logger.Info = false
	if err = chiselClient.Start(ctx); err != nil {
		return err
	}

	if err = chiselClient.Wait(); err != nil {
		return err
	}
	return nil
}
