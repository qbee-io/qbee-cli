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
	"testing"

	"go.qbee.io/client"
)

func TestDeviceInventory_IsMinimumAgentVersion(t *testing.T) {
	tests := []struct {
		name            string
		agentVersion    string
		requiredVersion string
		expected        bool
	}{
		// Agent version is greater
		{
			name:            "agent version major higher",
			agentVersion:    "2.0.0",
			requiredVersion: "1.0.0",
			expected:        true,
		},
		{
			name:            "agent version minor higher",
			agentVersion:    "1.5.0",
			requiredVersion: "1.4.0",
			expected:        true,
		},
		{
			name:            "agent version patch higher",
			agentVersion:    "1.0.5",
			requiredVersion: "1.0.3",
			expected:        true,
		},

		// Agent version is lower
		{
			name:            "agent version major lower",
			agentVersion:    "0.9.9",
			requiredVersion: "1.0.0",
			expected:        false,
		},
		{
			name:            "agent version minor lower",
			agentVersion:    "1.3.9",
			requiredVersion: "1.4.0",
			expected:        false,
		},
		{
			name:            "agent version patch lower",
			agentVersion:    "1.0.2",
			requiredVersion: "1.0.3",
			expected:        false,
		},

		// Exact match
		{
			name:            "versions are equal",
			agentVersion:    "1.2.3",
			requiredVersion: "1.2.3",
			expected:        true,
		},
		{
			name:            "single part versions equal",
			agentVersion:    "5",
			requiredVersion: "5",
			expected:        true,
		},

		// Agent has fewer parts
		{
			name:            "agent has fewer parts, required higher",
			agentVersion:    "1.0",
			requiredVersion: "1.0.5",
			expected:        false,
		},
		{
			name:            "agent has single part, required multi-part",
			agentVersion:    "1",
			requiredVersion: "1.0.1",
			expected:        false,
		},

		// Agent has more parts (should still work)
		{
			name:            "agent has more parts",
			agentVersion:    "1.0.0.5",
			requiredVersion: "1.0.0",
			expected:        true,
		},

		// Multiple part comparison
		{
			name:            "multi-part versions, agent higher",
			agentVersion:    "1.2.3.4",
			requiredVersion: "1.2.3.3",
			expected:        true,
		},
		{
			name:            "multi-part versions, agent lower",
			agentVersion:    "1.2.3.2",
			requiredVersion: "1.2.3.3",
			expected:        false,
		},

		// Invalid/non-numeric versions
		{
			name:            "invalid agent version",
			agentVersion:    "1.a.3",
			requiredVersion: "1.0.0",
			expected:        false,
		},
		{
			name:            "invalid required version",
			agentVersion:    "1.0.0",
			requiredVersion: "1.b.0",
			expected:        false,
		},
		{
			name:            "both invalid",
			agentVersion:    "a.b.c",
			requiredVersion: "x.y.z",
			expected:        false,
		},

		// Edge cases
		{
			name:            "empty agent version",
			agentVersion:    "",
			requiredVersion: "1.0.0",
			expected:        false,
		},
		{
			name:            "empty required version",
			agentVersion:    "1.0.0",
			requiredVersion: "",
			expected:        false,
		},
		{
			name:            "both empty",
			agentVersion:    "",
			requiredVersion: "",
			expected:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			di := client.DeviceInventory{
				SystemInfo: client.SystemInfo{
					AgentVersion: tt.agentVersion,
				},
			}

			result := di.IsMinimumAgentVersion(tt.requiredVersion)
			if result != tt.expected {
				t.Errorf("IsMinimumAgentVersion(%q) with agent version %q = %v, want %v",
					tt.requiredVersion, tt.agentVersion, result, tt.expected)
			}
		})
	}
}
