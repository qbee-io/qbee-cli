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

package client

import (
	"bytes"
	"encoding/json"
)

// UserBaseInfo is the base information of a user.
type UserBaseInfo struct {
	// ID is the unique identifier of the user.
	ID string `json:"user_id"`

	// FirstName is the first name of the user.
	FirstName string `json:"fname"`

	// LastName is the last name of the user.
	LastName string `json:"lname"`
}

// UnmarshalJSON implements custom unmarshalling to handle backward compatibility
// where UserBaseInfo could be represented as a string (user ID) or as an object.
func (u *UserBaseInfo) UnmarshalJSON(data []byte) error {
	// Check if data is a quoted string.
	if bytes.HasPrefix(data, []byte{0x22}) && bytes.HasSuffix(data, []byte{0x22}) {
		var userID string
		if err := json.Unmarshal(data, &userID); err != nil {
			return err
		}
		*u = UserBaseInfo{ID: userID}
		return nil
	}

	// Otherwise, unmarshal as the full object.
	type Alias UserBaseInfo

	var alias Alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	*u = UserBaseInfo(alias)

	return nil
}
