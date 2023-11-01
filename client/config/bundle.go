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

// Bundle defines a configuration bundle name.
type Bundle string

type BundleNames []Bundle

// Metadata for a configuration bundle.
type Metadata struct {
	// Enabled indicates whether the bundle is enabled.
	Enabled bool `json:"enabled"`

	// Extend indicates whether the bundle extends the configuration of its parent.
	Extend bool `json:"extend,omitempty"`

	// Version is the version of the bundle.
	// Currently, it's always set to v1.
	Version string `json:"version,omitempty"`

	// CommitID is the commit ID of the bundle.
	// This field is populated by the API.
	CommitID string `json:"bundle_commit_id,omitempty"`

	// Inherited indicates whether the bundle inherits from its parent.
	// DEPRECATED: This field is only used for backwards compatibility.
	Inherited bool `json:"inherited,omitempty"`

	// Inherits contains the complete inheritance information.
	// It starts with the node from which the bundle is inherited,
	// and ends with the node ID for which the configuration is rendered.
	// DEPRECATED: This field is only used for backwards compatibility.
	Inherits Ancestors `json:"inherits,omitempty"`

	// Reset instructs the system to reset the bundle configuration.
	Reset bool `json:"reset_to_group,omitempty"`
}

// BundleData combines all configuration bundles into one struct.
type BundleData struct {
	Settings             *Settings             `json:"settings,omitempty"`
	Users                *Users                `json:"users,omitempty"`
	SSHKeys              *SSHKeys              `json:"sshkeys,omitempty"`
	PackageManagement    *PackageManagement    `json:"package_management,omitempty"`
	FileDistribution     *FileDistribution     `json:"file_distribution,omitempty"`
	ConnectivityWatchdog *ConnectivityWatchdog `json:"connectivity_watchdog,omitempty"`
	ProcWatch            *ProcessWatch         `json:"proc_watch,omitempty"`
	NTP                  *NTP                  `json:"ntp,omitempty"`
	Parameters           *Parameters           `json:"parameters,omitempty"`
	SoftwareManagement   *SoftwareManagement   `json:"software_management,omitempty"`
	DockerContainers     *DockerContainers     `json:"docker_containers,omitempty"`
	Password             *Password             `json:"password,omitempty"`
	Firewall             *Firewall             `json:"firewall,omitempty"`
}
