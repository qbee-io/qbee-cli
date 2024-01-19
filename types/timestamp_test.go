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

package types

import (
	"encoding/json"
	"testing"
)

func TestTimestamp_UnmarshalJSON(t *testing.T) {
	type TestStruct struct {
		Value Timestamp `json:"t"`
	}

	tests := []struct {
		name     string
		json     string
		expected Timestamp
	}{
		{
			name:     "undefined",
			json:     `{}`,
			expected: Timestamp(0),
		},
		{
			name:     "from int",
			json:     `{"t": 123}`,
			expected: Timestamp(123),
		},
		{
			name:     "from unix timestamp string",
			json:     `{"t": "123"}`,
			expected: Timestamp(123),
		},
		{
			name:     "from date string",
			json:     `{"t": "Mon Jan  2 15:04:05 2006"}`,
			expected: Timestamp(1136214245),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TestStruct{}

			if err := json.Unmarshal([]byte(tt.json), &got); err != nil {
				t.Fatalf("UnmarshalJSON() error = %v", err)
			}

			if got.Value != tt.expected {
				t.Errorf("UnmarshalJSON() got = %v, want %v", got.Value, tt.expected)
			}
		})
	}
}
