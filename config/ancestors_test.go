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

package config

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestAncestorsChain_MarshalJSON(t *testing.T) {
	data := struct {
		Ancestors Ancestors `json:"ancestors,omitempty"`
	}{
		Ancestors: Ancestors{
			{
				NodeID: "root",
				Title:  "Root",
			},
			{
				NodeID: "group",
				Title:  "Group",
			},
		},
	}

	got, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `{"ancestors":{"root":"Root","group":"Group"}}`

	if !reflect.DeepEqual(string(got), expected) {
		t.Fatalf("expected: %s, got: %s", expected, string(got))
	}
}

func TestAncestorsChain_UnmarshalJSON(t *testing.T) {
	expectedData := struct {
		Ancestors Ancestors `json:"ancestors,omitempty"`
	}{
		Ancestors{
			{
				NodeID: "root",
				Title:  `Root "test"`,
			},
			{
				NodeID: "group",
				Title:  "Group, All",
			},
		},
	}

	unmarshalled := struct {
		Ancestors Ancestors `json:"ancestors,omitempty"`
	}{}

	marshalledWithWhitespaces, err := json.MarshalIndent(expectedData, " ", " ")
	if err != nil {
		t.Fatalf("error marshalling: %v", err)
	}

	if err = json.Unmarshal(marshalledWithWhitespaces, &unmarshalled); err != nil {
		t.Fatalf("error unmarshaling: %v", err)
	}

	if !reflect.DeepEqual(unmarshalled, expectedData) {
		t.Fatalf("expected: %s, got: %s", expectedData, unmarshalled)
	}
}

func TestAncestorsChain_MarshalJSON_Empty(t *testing.T) {
	data := struct {
		Ancestors Ancestors `json:"ancestors,omitempty"`
	}{}

	got, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `{}`

	if string(got) != expected {
		t.Fatalf("expected: %s, got: %s", expected, string(got))
	}

	unmarshalled := struct {
		Ancestors Ancestors `json:"ancestors,omitempty"`
	}{}

	if err = json.Unmarshal(got, &unmarshalled); err != nil {
		t.Fatalf("error unmarshaling: %v", err)
	}

	if !reflect.DeepEqual(unmarshalled, data) {
		t.Fatalf("expected: %s, got: %s", data, unmarshalled)
	}
}
