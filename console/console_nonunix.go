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

//go:build !unix

package console

import (
	"context"

	"github.com/xtaci/smux"
)

// ResizeConsole resizes the terminal. On non-Unix systems this is a no-op.
func ResizeConsole(ctx context.Context, shellStream *smux.Session, sessionID string, termFD, width, height int) {
	return
}
