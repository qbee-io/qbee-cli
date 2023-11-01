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

package client_test

import (
	"reflect"
	"testing"

	"github.com/qbee-io/qbee-cli/client"
)

// TestParseRemoteAccessTarget tests that we return meaningful errors for invalid targets.
func TestParseRemoteAccessTarget(t *testing.T) {
	tests := []struct {
		name         string
		targetString string
		want         client.RemoteAccessTarget
		wantErr      string
	}{
		{
			name:         "valid tcp target",
			targetString: "tcp:1:af973d5836408b20cf051342f1cbf75a1f9096385993f15a4645e4f75e75f288:2",
			want: client.RemoteAccessTarget{
				Protocol:   "tcp",
				LocalPort:  1,
				Device:     "af973d5836408b20cf051342f1cbf75a1f9096385993f15a4645e4f75e75f288",
				RemotePort: 2,
			},
		},
		{
			name:         "valid udp target",
			targetString: "udp:1:af973d5836408b20cf051342f1cbf75a1f9096385993f15a4645e4f75e75f288:2",
			want: client.RemoteAccessTarget{
				Protocol:   "udp",
				LocalPort:  1,
				Device:     "af973d5836408b20cf051342f1cbf75a1f9096385993f15a4645e4f75e75f288",
				RemotePort: 2,
			},
		},
		{
			name:         "invalid format",
			targetString: "af973d5836408b20cf051342f1cbf75a1f9096385993f15a4645e4f75e75f288",
			wantErr:      "invalid format",
		},
		{
			name:         "unsupported protocol",
			targetString: "x:1:af973d5836408b20cf051342f1cbf75a1f9096385993f15a4645e4f75e75f288:2",
			wantErr:      "invalid protocol",
		},
		{
			name:         "local port out of range",
			targetString: "tcp:123456:af973d5836408b20cf051342f1cbf75a1f9096385993f15a4645e4f75e75f288:2",
			wantErr:      "invalid local port: invalid port number",
		},
		{
			name:         "local port does not support service name",
			targetString: "tcp:ssh:af973d5836408b20cf051342f1cbf75a1f9096385993f15a4645e4f75e75f288:2",
			wantErr:      "invalid local port: invalid port number",
		},
		{
			name:         "empty local port",
			targetString: "tcp::af973d5836408b20cf051342f1cbf75a1f9096385993f15a4645e4f75e75f288:2",
			wantErr:      "invalid local port: empty port",
		},
		{
			name:         "remote port out of range",
			targetString: "tcp:1:af973d5836408b20cf051342f1cbf75a1f9096385993f15a4645e4f75e75f288:234567",
			wantErr:      "invalid remote port: invalid port number",
		},
		{
			name:         "remote port does not support service name",
			targetString: "tcp:1:af973d5836408b20cf051342f1cbf75a1f9096385993f15a4645e4f75e75f288:ssh",
			wantErr:      "invalid remote port: invalid port number",
		},
		{
			name:         "empty remote port",
			targetString: "tcp:1:af973d5836408b20cf051342f1cbf75a1f9096385993f15a4645e4f75e75f288:",
			wantErr:      "invalid remote port: empty port",
		},
		{
			name:         "invalid device ID length",
			targetString: "tcp:1:abc:2",
			wantErr:      "invalid device",
		},
		{
			name:         "invalid device characters",
			targetString: "tcp:1:xx973d5836408b20cf051342f1cbf75a1f9096385993f15a4645e4f75e75f288:2",
			wantErr:      "invalid device",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.ParseRemoteAccessTarget(tt.targetString)
			if err != nil && err.Error() != tt.wantErr {
				t.Errorf("ParseRemoteAccessTarget() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && tt.wantErr != "" {
				t.Errorf("ParseRemoteAccessTarget() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRemoteAccessTarget() got = %v, want %v", got, tt.want)
			}
		})
	}
}
