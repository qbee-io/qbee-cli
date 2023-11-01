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

const FileDistributionBundle Bundle = "file_distribution"

// FileDistribution controls files in the system.
//
// Example payload:
//
//	{
//	 "files": [
//	   {
//	     "pre_condition": "/tmp/test.sh",
//	     "templates": [
//	       {
//	         "source": "demo_file.json",
//	         "destination": "/tmp/demo_file.json",
//	         "is_template": true
//	       }
//	     ],
//	     "parameters": [
//	       {
//	         "key": "VAR1",
//	         "value": "VAL1"
//	       }
//	     ],
//	     "command": "echo \"it worked!\""
//	   }
//	 ]
//	}
type FileDistribution struct {
	Metadata

	FileSets []FileSet `json:"files,omitempty"`
}

// FileSet defines a file set to be maintained in the system.
type FileSet struct {
	// Files defines files to be created in the filesystem.
	Files []File `json:"templates,omitempty"`

	// Parameters define values to be used for template files.
	TemplateParameters []TemplateParameter `json:"parameters,omitempty"`

	// AfterCommand defines a command to be executed after files are saved on the filesystem.
	AfterCommand string `json:"command,omitempty"`

	// PreCondition defines an optional command which needs to return 0 in order for the FileSet to be executed.
	PreCondition string `json:"pre_condition,omitempty"`
}

// File defines a single file parameters.
type File struct {
	// Source full file path from the file manager.
	Source string `json:"source,omitempty"`

	// Destination defines absolute path of the file in the filesystem.
	Destination string `json:"destination,omitempty"`

	// IsTemplate defines whether the file should be processed by the templating engine.
	IsTemplate bool `json:"is_template"`
}

// TemplateParameter defines a single parameter used to replace placeholder in a template.
type TemplateParameter struct {
	// Key of the parameter used in files.
	Key string `json:"key,omitempty"`

	// Value of the parameter which will replace Key placeholders.
	Value string `json:"value,omitempty"`
}
