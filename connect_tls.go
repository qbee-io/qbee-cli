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

//go:build !insecure_edge

package client

import "crypto/tls"

// edgeTLSConfig returns the TLS client configuration for edge connections.
//
// Production builds always enforce certificate verification, so this returns
// nil and lets the transport use its secure defaults. A development-only
// variant that skips verification for local edges can be enabled with the
// "insecure_edge" build tag (see connect_tls_insecure.go).
func edgeTLSConfig(_ string) *tls.Config {
	return nil
}
