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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"go.qbee.io/transport"

	chisel "github.com/jpillora/chisel/client"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/http/httpproxy"
	"golang.org/x/term"
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

// connectStdio connects to the given target using stdin/stdout.
func (cli *Client) connectStdio(ctx context.Context, client *transport.Client, target RemoteAccessTarget) error {
	remoteHostPort := fmt.Sprintf("%s:%s", target.RemoteHost, target.RemotePort)

	stream, err := client.OpenStream(ctx, transport.MessageTypeTCPTunnel, []byte(remoteHostPort))
	if err != nil {
		return fmt.Errorf("error opening stream: %w", err)
	}
	defer stream.Close()

	// copy from stdin to stream
	go func() {
		_, _ = io.Copy(stream, os.Stdin)
	}()

	// copy from stream to stdout
	_, err = io.Copy(os.Stdout, stream)

	return err
}

// connect establishes a connection to a remote device.
func (cli *Client) connect(ctx context.Context, deviceUUID, edgeHost string, targets []RemoteAccessTarget) error {
	edgeURL := fmt.Sprintf("https://%s/device/%s", edgeHost, deviceUUID)

	var tlsConfig *tls.Config

	// for testing purposes, allow connections to localhost without verifying the certificate
	if strings.HasPrefix(edgeHost, "edge:") || strings.HasPrefix(edgeHost, "localhost:") {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	if len(targets) == 0 {
		return fmt.Errorf("no targets defined")
	}

	client, err := transport.NewClient(ctx, edgeURL, cli.authToken, tlsConfig)
	if err != nil {
		return fmt.Errorf("error initializing remote access client: %w", err)
	}

	// close the client and all local listeners when the context is cancelled
	closers := []io.Closer{client}
	defer func() {
		for _, closer := range closers {
			_ = closer.Close()
		}
	}()

	if len(targets) == 1 && targets[0].LocalPort == "stdio" {
		return cli.connectStdio(ctx, client, targets[0])
	}

	for _, target := range targets {
		if target.LocalPort == "stdio" {
			return fmt.Errorf("stdio is only supported for single target connections")
		}

		localHostPort := fmt.Sprintf("localhost:%s", target.LocalPort)
		remoteHostPort := fmt.Sprintf("%s:%s", target.RemoteHost, target.RemotePort)

		switch target.Protocol {
		case "tcp":
			var tcpListener *net.TCPListener
			if tcpListener, err = client.OpenTCPTunnel(ctx, localHostPort, remoteHostPort); err != nil {
				return fmt.Errorf("error opening TCP tunnel: %w", err)
			}

			closers = append(closers, tcpListener)
		case "udp":
			var udpConn *transport.UDPTunnel
			if udpConn, err = client.OpenUDPTunnel(ctx, localHostPort, remoteHostPort); err != nil {
				return fmt.Errorf("error opening UDP tunnel: %w", err)
			}

			closers = append(closers, udpConn)
		default:
			return fmt.Errorf("invalid protocol %s", target.Protocol)
		}

		fmt.Printf("Tunneling %s %s to %s\n", target.Protocol, localHostPort, remoteHostPort)
	}

	// Wait for context to be cancelled
	<-ctx.Done()

	return nil
}

// legacyConnect establishes a connection to a remote device using the legacy remote access solution.
func (cli *Client) legacyConnect(ctx context.Context, deviceID string, targets []RemoteAccessTarget) error {
	ports := make([]string, len(targets))
	for _, target := range targets {
		// only localhost is supported as remote host for legacy remote access
		if target.RemoteHost != "localhost" {
			return fmt.Errorf("invalid remote host: only localhost is supported")
		}

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

// Connect establishes a connection to a remote device.
func (cli *Client) Connect(ctx context.Context, deviceID string, targets []RemoteAccessTarget) error {
	deviceStatus, err := cli.GetDeviceStatus(ctx, deviceID)
	if err != nil {
		return err
	}

	if !deviceStatus.RemoteAccess {
		return fmt.Errorf("remote access is not available for device %s", deviceID)
	}

	switch deviceStatus.EdgeVersion {
	case EdgeVersionOpenVPN:
		return cli.legacyConnect(ctx, deviceID, targets)
	case EdgeVersionNative:
		return cli.connect(ctx, deviceStatus.UUID, deviceStatus.Edge, targets)
	default:
		return fmt.Errorf("unsupported edge version %d", deviceStatus.EdgeVersion)
	}
}

// ConnectShell establishes a shell connection to a remote device.
func (cli *Client) ConnectShell(ctx context.Context, deviceID string) error {

	deviceStatus, err := cli.GetDeviceStatus(ctx, deviceID)
	if err != nil {
		return err
	}

	if !deviceStatus.RemoteAccess {
		return fmt.Errorf("remote access is not available for device %s", deviceID)
	}

	edgeURL := fmt.Sprintf("https://%s/device/%s", deviceStatus.Edge, deviceStatus.UUID)

	var tlsConfig *tls.Config

	client, err := transport.NewClient(ctx, edgeURL, cli.authToken, tlsConfig)
	if err != nil {
		return fmt.Errorf("error initializing remote access client: %w", err)
	}

	// close the client and all local listeners when the context is cancelled
	closers := []io.Closer{client}
	defer func() {
		for _, closer := range closers {
			_ = closer.Close()
		}
	}()

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("terminal make raw: %s", err)
	}
	defer term.Restore(fd, oldState)

	w, h, err := terminal.GetSize(fd)
	if err != nil {
		return fmt.Errorf("terminal get size: %s", err)
	}

	var initCmd = &transport.PTYCommand{
		Type:      transport.PTYCommandTypeResize,
		SessionID: "",
		Cols:      uint16(w),
		Rows:      uint16(h),
	}

	payload, err := json.Marshal(initCmd)
	if err != nil {
		return fmt.Errorf("error marshaling initial window size: %w", err)
	}

	shellStream, err := client.OpenStream(ctx, transport.MessageTypePTY, payload)
	if err != nil {
		return fmt.Errorf("error opening shell stream: %w", err)
	}
	defer shellStream.Close()

	// copy from stdin to stream

	stdOutClosed := make(chan bool)
	stdInClosed := make(chan bool)

	go func() {
		var buf [1024]byte
		for {
			n, err := os.Stdin.Read(buf[:])
			if err != nil {
				fmt.Printf("error reading from stdin: %s\n", err)
				stdInClosed <- true
				return
			}
			_, err = shellStream.Write(buf[:n])

			if err != nil {
				if err == io.EOF {
					stdInClosed <- true
					return
				}
				fmt.Printf("error writing to stream: %s\n", err)
				stdInClosed <- true
				return
			}
		}
	}()

	go func() {
		var buf [1024]byte
		for {
			n, err := shellStream.Read(buf[:])
			if err != nil {
				if err == io.EOF {
					stdOutClosed <- true
					return
				}
				fmt.Printf("error reading from stream: %s\n", err)
				stdOutClosed <- true
				return
			}
			_, err = os.Stdout.Write(buf[:n])
			if err != nil {
				fmt.Printf("error writing to stdout: %s\n", err)
				stdOutClosed <- true
				return
			}
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-stdOutClosed:
			return nil
		case <-stdInClosed:
			return nil
		}
	}
}
