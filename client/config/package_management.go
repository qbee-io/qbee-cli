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

// PackageManagementBundle defines name for the package management bundle.
const PackageManagementBundle Bundle = "package_management"

// PackageManagement controls system packages.
//
// Example payload:
//
//	{
//	 "pre_condition": "test command",
//	 "items": [
//	   {
//	     "name": "httpd2",
//	     "version": "1.2.3"
//	   }
//	 ],
//	 "reboot_mode": "always",
//	 "full_upgrade": false
//	}
type PackageManagement struct {
	Metadata

	// PreCondition if set will be executed before package maintenance.
	// If the command returns a non-zero exit code, the package maintenance will be skipped.
	PreCondition string `json:"pre_condition,omitempty"`

	// RebootMode defines whether system should be rebooted after package maintenance or not.
	RebootMode RebootMode `json:"reboot_mode"`

	// FullUpgrade if set to true will perform a full system upgrade.
	FullUpgrade bool `json:"full_upgrade"`

	// Packages defines a list of packages to be maintained.
	Packages []Package `json:"items,omitempty"`
}

// RebootMode defines whether system should be rebooted after package maintenance or not.
type RebootMode string

// Supported reboot modes.
const (
	RebootNever  RebootMode = "never"
	RebootAlways RebootMode = "always"
)

// Package defines a package to be maintained in the system.
type Package struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}
