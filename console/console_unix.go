// Copyright 2024 qbee.io
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

//go:build unix

package console

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/xtaci/smux"
	"go.qbee.io/transport"
	"golang.org/x/term"
)

// ResizeConsole resizes the terminal. On non-Unix systems this is a no-op.
func ResizeConsole(ctx context.Context, session *smux.Session, sessionID string, termFD, width, height int) {

	windowChange := make(chan os.Signal, 1)
	signal.Notify(windowChange, syscall.SIGWINCH)
	defer signal.Stop(windowChange)

	for {
		select {
		case <-ctx.Done():
			return
		case <-windowChange:
			newWidth, newHeight, _ := term.GetSize(termFD)

			if newWidth == width && newHeight == height {
				return
			}

			width = newWidth
			height = newHeight

			if err := sendResizeCommand(session, sessionID, width, height); err != nil {
				fmt.Fprintf(os.Stderr, "error resizing window: %v\n", err)
				return
			}
		}
	}
}

// sendRezizeCommand sends a resize command to the remote shell.
func sendResizeCommand(session *smux.Session, sessionID string, width, height int) error {
	cmd := transport.PTYCommand{
		Type:      transport.PTYCommandTypeResize,
		SessionID: sessionID,
		Cols:      uint16(width),
		Rows:      uint16(height),
	}

	shellStream, err := session.OpenStream()
	if err != nil {
		return fmt.Errorf("error opening shell stream: %w", err)
	}

	var payload []byte
	if payload, err = json.Marshal(cmd); err != nil {
		return fmt.Errorf("error marshaling window resize command: %w", err)
	}

	if err := transport.WriteMessage(shellStream, transport.MessageTypePTYCommand, payload); err != nil {
		return fmt.Errorf("error writing window resize command: %w", err)
	}

	if _, err = transport.ExpectOK(shellStream); err != nil {
		return fmt.Errorf("error resizing window: %w", err)
	}
	return nil
}
