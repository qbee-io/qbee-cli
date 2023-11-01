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
	"bytes"
	"encoding/json"
	"sort"
)

// Ancestor defines a node in the Ancestors.
type Ancestor struct {
	NodeID string
	Title  string
}

// Ancestors represents a list of ancestors (root is first).
type Ancestors []Ancestor

// MarshalJSON allows to marshal ancestors chain as an ordered map.
//
//goland:noinspection GoMixedReceiverTypes
func (ancestors Ancestors) MarshalJSON() ([]byte, error) {
	buf := new(bytes.Buffer)

	buf.WriteByte('{')

	for i, ancestor := range ancestors {
		if i > 0 {
			buf.WriteByte(',')
		}

		if nodeID, err := json.Marshal(ancestor.NodeID); err != nil {
			return nil, err
		} else {
			buf.Write(nodeID)
		}

		buf.WriteByte(':')

		if title, err := json.Marshal(ancestor.Title); err != nil {
			return nil, err
		} else {
			buf.Write(title)
		}
	}

	buf.WriteByte('}')

	return buf.Bytes(), nil
}

// UnmarshalJSON allows to unmarshal ancestors chain from an ordered map.
//
//goland:noinspection GoMixedReceiverTypes
func (ancestors *Ancestors) UnmarshalJSON(data []byte) error {
	// decode as a map first, so JSON module does all the input validation for us
	dataMap := make(map[string]string)
	if err := json.Unmarshal(data, &dataMap); err != nil {
		return err
	}

	// for empty map, return empty slice
	if len(dataMap) == 0 {
		return nil
	}

	// to order the ancestors, lookup indexes of keys in the byte-string and use those as sorting keys
	sortingSlice := make([]struct {
		ancestor Ancestor
		index    int
	}, 0, len(dataMap))

	for k, v := range dataMap {
		rawKey, err := json.Marshal(k)
		if err != nil {
			return err
		}

		ancestor := struct {
			ancestor Ancestor
			index    int
		}{
			ancestor: Ancestor{
				NodeID: k,
				Title:  v,
			},
			index: bytes.Index(data, rawKey),
		}

		sortingSlice = append(sortingSlice, ancestor)
	}

	sort.Slice(sortingSlice, func(i, j int) bool {
		return sortingSlice[i].index < sortingSlice[j].index
	})

	// apply ancestors in the right order
	for i := range sortingSlice {
		*ancestors = append(*ancestors, sortingSlice[i].ancestor)
	}

	return nil
}
