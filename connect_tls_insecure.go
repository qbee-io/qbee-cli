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

//go:build insecure_edge

package client

import (
	"crypto/tls"
	"strings"
)

// edgeTLSConfig returns the TLS client configuration for edge connections.
//
// WARNING: this variant is compiled only with the "insecure_edge" build tag.
// It disables TLS certificate verification for local development edges (the
// docker-compose "edge:" service and "localhost:"). Binaries built with this
// tag are vulnerable to man-in-the-middle attacks and must never be released.
func edgeTLSConfig(edgeHost string) *tls.Config {
	if strings.HasPrefix(edgeHost, "edge:") || strings.HasPrefix(edgeHost, "localhost:") {
		return &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec // dev-only, guarded by the insecure_edge build tag
		}
	}

	return nil
}
