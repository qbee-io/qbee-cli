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
			targetString: "1:localhost:2",
			want: client.RemoteAccessTarget{
				Protocol:   "tcp",
				LocalPort:  "1",
				RemoteHost: "localhost",
				RemotePort: "2",
			},
		},
		{
			name:         "valid udp target",
			targetString: "1:localhost:2/udp",
			want: client.RemoteAccessTarget{
				Protocol:   "udp",
				LocalPort:  "1",
				RemoteHost: "localhost",
				RemotePort: "2",
			},
		},
		{
			name:         "valid tcp target",
			targetString: "stdio:localhost:2",
			want: client.RemoteAccessTarget{
				Protocol:   "tcp",
				LocalPort:  "stdio",
				RemoteHost: "localhost",
				RemotePort: "2",
			},
		},
		{
			name:         "valid udp target",
			targetString: "stdio:localhost:2/udp",
			want: client.RemoteAccessTarget{
				Protocol:   "udp",
				LocalPort:  "stdio",
				RemoteHost: "localhost",
				RemotePort: "2",
			},
		},
		{
			name:         "invalid format",
			targetString: "localhost",
			wantErr:      "invalid format",
		},
		{
			name:         "unsupported host",
			targetString: "123:example.com:123",
			wantErr:      "invalid remote host: only localhost is supported",
		},
		{
			name:         "local port out of range",
			targetString: "123456:localhost:2",
			wantErr:      "invalid local port: invalid port number",
		},
		{
			name:         "local port does not support service name",
			targetString: "ssh:localhost:2",
			wantErr:      "invalid local port: invalid port number",
		},
		{
			name:         "empty local port",
			targetString: ":localhost:2",
			wantErr:      "invalid local port: empty port",
		},
		{
			name:         "remote port out of range",
			targetString: "1:localhost:234567",
			wantErr:      "invalid remote port: invalid port number",
		},
		{
			name:         "remote port does not support service name",
			targetString: "1:localhost:ssh",
			wantErr:      "invalid remote port: invalid port number",
		},
		{
			name:         "empty remote port",
			targetString: "1:localhost:",
			wantErr:      "invalid remote port: empty port",
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
