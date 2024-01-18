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

// SoftwareManagementBundle defines name for the software management bundle.
const SoftwareManagementBundle Bundle = "software_management"

// SoftwareManagement controls software in the system.
//
// Example payload:
//
//	{
//	 "items": [
//	   {
//	     "package": "pkg1",
//	     "service_name": "serviceName",
//	     "config_files": [
//	       {
//	         "config_template": "configFileTemplate",
//	         "config_location": "configFileLocation"
//	       }
//	     ],
//	     "parameters": [
//	       {
//	         "key": "configKey",
//	         "value": "configValue"
//	       }
//	     ]
//	   }
//	 ]
//	}
type SoftwareManagement struct {
	Metadata

	// Items to be installed.
	Items []SoftwarePackage `json:"items,omitempty"`
}

// SoftwarePackage defines software package to be maintained in the system.
type SoftwarePackage struct {
	// Package defines a package name to install.
	Package string `json:"package,omitempty"`

	// ServiceName defines an optional service name (if empty, Package is used).
	ServiceName string `json:"service_name,omitempty"`

	// PreCondition defines an optional command which needs to return 0 in order for the SoftwarePackage to be installed.
	PreCondition string `json:"pre_condition,omitempty"`

	// ConfigFiles to be created for the software.
	ConfigFiles []ConfigurationFile `json:"config_files,omitempty"`

	// Parameters for the ConfigFiles templating.
	Parameters []ConfigurationFileParameter `json:"parameters,omitempty"`
}

// ConfigurationFile definition.
type ConfigurationFile struct {
	// ConfigTemplate defines a source template file from file manager.
	ConfigTemplate string `json:"config_template,omitempty"`

	// ConfigLocation defines an absolute path in the system where file will be created.
	ConfigLocation string `json:"config_location,omitempty"`
}

// ConfigurationFileParameter defines parameter to be used in ConfigurationFile.
type ConfigurationFileParameter struct {
	// Key defines parameters name.
	Key string `json:"key,omitempty"`

	// Value defines parameters value.
	Value string `json:"value,omitempty"`
}
