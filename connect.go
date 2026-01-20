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
	"math"
	"math/rand/v2"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"go.qbee.io/client/console"
	"go.qbee.io/transport"

	"golang.org/x/term"
)

// ConnectMulti establishes connections to multiple remote devices concurrently.
func (cli *Client) ConnectMulti(ctx context.Context, connections []RemoteAccessConnection, allowFailures bool) error {
	return cli.ConnectMultiRetry(ctx, connections, allowFailures, 1)
}

// ConnectMultiRetry establishes connections to multiple remote devices concurrently with retries.
func (cli *Client) ConnectMultiRetry(ctx context.Context, connections []RemoteAccessConnection, allowFailures bool, retries int) error {
	wg := sync.WaitGroup{}
	errChan := make(chan error)
	done := make(chan bool)

	for _, conn := range connections {
		wg.Add(1)

		go func(connection RemoteAccessConnection) {
			defer wg.Done()

			if err := cli.ParseConnectRetry(ctx, connection.DeviceID, connection.Targets, retries); err != nil {
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
	return cli.ParseConnectRetry(ctx, deviceID, targets, 1)
}

// ParseConnectRetry parses a device ID and a list of targets and establishes a connection to the device with retries.
func (cli *Client) ParseConnectRetry(ctx context.Context, deviceID string, targets []string, retries int) error {

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

	if retries < 0 {
		return fmt.Errorf("retries must be a positive number")
	}

	var err error
	attempts := 0
	baseTime := 5 * time.Second
	maxBackoff := 1 * time.Minute
	for {
		err = cli.Connect(ctx, deviceID, parsedTargets)
		backoff := time.Duration(math.Min(float64(baseTime)*math.Pow(2, float64(attempts)), float64(maxBackoff)))

		if err != nil {
			fmt.Printf("error connecting to device %s: %s\n", deviceID, err)
		}

		attempts++
		// Exit if the maximum number of retries has been reached.
		if attempts >= retries && retries > 0 {
			break
		}

		jitter := time.Duration(rand.Float64() * float64(backoff) * 1)
		nextBackoff := backoff + jitter
		fmt.Printf("Attempt %d failed. Retrying in %v...\n", attempts, nextBackoff)
		time.Sleep(nextBackoff)

	}

	return err
}

// connectStdio connects to the given target using stdin/stdout.
func (cli *Client) connectStdio(ctx context.Context, client *transport.Client, target RemoteAccessTarget) error {
	remoteHostPort := fmt.Sprintf("%s:%s", target.RemoteHost, target.RemotePort)

	stream, err := client.OpenStream(ctx, transport.MessageTypeTCPTunnel, []byte(remoteHostPort))
	if err != nil {
		return fmt.Errorf("error opening stream: %w", err)
	}
	defer func() { _ = stream.Close() }()

	// copy from stdin to stream
	go func() {
		_, _ = io.Copy(stream, os.Stdin)
	}()

	// copy from stream to stdout
	_, err = io.Copy(os.Stdout, stream)

	return err
}

// getConnectClient gets a transport client for the device connection
func (cli *Client) getConnectClient(ctx context.Context, deviceID string) (*transport.Client, error) {
	deviceStatus, err := cli.GetDeviceStatus(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	if !deviceStatus.RemoteAccess {
		return nil, fmt.Errorf("remote access is not available for device %s", deviceID)
	}

	edgeURL := fmt.Sprintf("https://%s/device/%s", deviceStatus.Edge, deviceStatus.UUID)

	var tlsConfig *tls.Config

	if strings.HasPrefix(deviceStatus.Edge, "edge:") || strings.HasPrefix(deviceStatus.Edge, "localhost:") {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	client, err := transport.NewClient(ctx, edgeURL, cli.authToken, tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("error initializing remote access client: %w", err)
	}

	return client, nil
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

		localHostPort := fmt.Sprintf("%s:%s", target.LocalHost, target.LocalPort)
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

	smuxSession := client.GetSession()
	errChan := make(chan error)

	// block until the session is closed or an error occurs. Typically this will happen when
	// the device is disconnected mid-session.
	go func() {
		for {
			_, err := smuxSession.AcceptStream()
			if err != nil {
				errChan <- err
				return
			}
		}
	}()

	select {
	case err := <-errChan:
		_ = smuxSession.Close()
		return fmt.Errorf("session error for device %s: %w", deviceUUID, err)
	case <-ctx.Done():
		fmt.Printf("Connection closed\n")
		return nil
	}
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
	case EdgeVersionNative:
		return cli.connect(ctx, deviceStatus.UUID, deviceStatus.Edge, targets)
	default:
		return fmt.Errorf("unsupported edge version %d", deviceStatus.EdgeVersion)
	}
}

// ConnectTerminal establishes a shell connection to a remote device.
func (cli *Client) ConnectTerminal(ctx context.Context, deviceID string, command []string) error {

	termStdinFd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(termStdinFd)
	if err != nil {
		return fmt.Errorf("terminal make raw: %s", err)
	}

	defer func() {
		err := term.Restore(termStdinFd, oldState)
		if err != nil {
			fmt.Printf("error restoring terminal state: %s\n", err)
		}
	}()

	termWidth, termHeight, err := term.GetSize(termStdinFd)
	if err != nil {
		return fmt.Errorf("terminal get size: %s", err)
	}

	var initCmd = &transport.PTYCommand{
		Type:      transport.PTYCommandTypeResize,
		SessionID: "",
		Cols:      uint16(termWidth),
		Rows:      uint16(termHeight),
	}

	if command != nil {
		initCmd.Command = command[0]

		if len(command) > 1 {
			initCmd.CommandArgs = command[1:]
		}
	}

	client, err := cli.getConnectClient(ctx, deviceID)
	if err != nil {
		return err
	}

	// close the client and all local listeners when the context is cancelled
	closers := []io.Closer{client}
	defer func() {
		for _, closer := range closers {
			_ = closer.Close()
		}
	}()

	payload, err := json.Marshal(initCmd)
	if err != nil {
		return fmt.Errorf("error marshaling initial window size: %w", err)
	}

	shellStream, sessionIDBytes, err := client.OpenStreamPayload(ctx, transport.MessageTypePTY, payload)

	if err != nil {
		return fmt.Errorf("error opening shell stream: %w", err)
	}

	defer func() { _ = shellStream.Close() }()

	go console.ResizeConsole(ctx, client.GetSession(), string(sessionIDBytes), termStdinFd, termWidth, termHeight)

	errChan := make(chan error)

	go readerLoop(shellStream, os.Stdout, errChan)
	go readerLoop(os.Stdin, shellStream, errChan)

	select {
	case <-ctx.Done():
		return nil
	case err := <-errChan:
		if err != nil {
			return err
		}
		return nil
	}
}

// readerLoop reads from reader and writes to writer until EOF or an error occurs.
func readerLoop(in io.Reader, out io.Writer, errChan chan error) {
	var buf [1024]byte

	for {
		n, err := in.Read(buf[:])
		if err != nil {
			if err == io.EOF {
				errChan <- nil
				return
			}
			errChan <- err
			return
		}

		_, err = out.Write(buf[:n])
		if err != nil {
			errChan <- err
			return
		}

	}
}
