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
	"bytes"
	"encoding/json"
	"strconv"
	"time"
)

// Timestamp is an int64-based type that allows to unmarshal also from string containing unix timestamp or date string.
type Timestamp int64

func parseTimestamp(s string) (Timestamp, error) {
	if s == "" {
		return 0, nil
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		timeValue, timeErr := time.Parse(time.ANSIC, s)
		if timeErr != nil {
			return 0, err
		}

		return Timestamp(timeValue.Unix()), nil
	}

	return Timestamp(i), nil
}

// UnmarshalJSON allows to unmarshal Int64String from string or int64.
func (i64 *Timestamp) UnmarshalJSON(data []byte) error {
	if bytes.HasPrefix(data, []byte(`"`)) {
		var value string
		var err error

		if err = json.Unmarshal(data, &value); err != nil {
			return err
		}

		*i64, err = parseTimestamp(value)

		return err
	}

	var value int64
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	*i64 = Timestamp(value)

	return nil
}

// Time returns the time.Time value of the timestamp.
func (i64 *Timestamp) Time() time.Time {
	return time.Unix(int64(*i64), 0)
}
