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

// SettingsBundle defines name for the agent settings bundle.
const SettingsBundle Bundle = "settings"

// Settings
//
// Example payload:
//
//	"settings": {
//	  "metrics": true,
//	  "reports": true,
//	  "remoteconsole": true,
//	  "software_inventory": true,
//	  "process_inventory": true,
//	  "agentinterval": 10
//	}
type Settings struct {
	Metadata

	// Metrics collection enabled.
	Metrics bool `json:"metrics"`

	// Reports collection enabled.
	Reports bool `json:"reports"`

	// RemoteConsole access enabled.
	RemoteConsole bool `json:"remoteconsole"`

	// SoftwareInventory collection enabled.
	SoftwareInventory bool `json:"software_inventory"`

	// ProcessInventory collection enabled.
	ProcessInventory bool `json:"process_inventory"`

	// AgentInterval defines how often agent reports back to the device hub (in minutes).
	AgentInterval int `json:"agentinterval"`
}
