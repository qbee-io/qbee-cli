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
func ResizeConsole(ctx context.Context, shellStream *smux.Stream, sessionID string, termFD, width, height int) {

	windowChange := make(chan os.Signal, 1)
	signal.Notify(windowChange, syscall.SIGWINCH)
	defer signal.Stop(windowChange)

	for {
		select {
		case <-ctx.Done():
			return
		case <-windowChange:
			newWidth, newHeight, _ := term.GetSize(termFD)
			if newWidth != width || newHeight != height {
				width = newWidth
				height = newHeight

				cmd := transport.PTYCommand{
					Type:      transport.PTYCommandTypeResize,
					SessionID: sessionID,
					Cols:      uint16(width),
					Rows:      uint16(height),
				}
				var payload []byte
				var err error
				if payload, err = json.Marshal(cmd); err != nil {
					fmt.Printf("error marshaling window resize command: %s\n", err)
					return
				}

				if err := transport.WriteMessage(shellStream, transport.MessageTypePTY, payload); err != nil {
					fmt.Printf("error writing window resize command: %s\n", err)
					return
				}

				if _, err = transport.ExpectOK(shellStream); err != nil {
					fmt.Printf("error resizing window: %s\n", err)
					return
				}
			}
		}
	}
}
