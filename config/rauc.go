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

package config

// RaucBundle defines the name for the RAUC bundle.
const RaucBundle Bundle = "rauc"

// Rauc configures an A/B system update using RAUC.
//
// example payload
//
//	{
//	  "pre_condition": "true",
//	  "rauc_bundle": "/path/to/bundle.raucb",
//	}
type Rauc struct {
	Metadata

	// PreCondition defines an optional command which needs to return 0 in order for RAUC bundle to be installed.
	PreCondition string `json:"pre_condition,omitempty"`

	// RaucBundle defines the rauc bundle to be installed.
	RaucBundle string `json:"rauc_bundle"`
}
