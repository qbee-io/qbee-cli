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

package client

import (
	"encoding/json"
	"fmt"
)

type Error map[string]any

func (error Error) Error() string {
	data, err := json.Marshal(error)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s", data)
}

// ParseErrorResponse parses API error from the provided response body.
func ParseErrorResponse(responseBody []byte) Error {
	apiError := make(Error)

	if err := json.Unmarshal(responseBody, &apiError); err != nil {
		panic(fmt.Errorf("failed to decode error from (%v): %s", err, responseBody))
	}

	return apiError
}
