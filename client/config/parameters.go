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

const ParametersBundle Bundle = "parameters"

// Parameters sets global configuration parameters.
//
// Example payload:
//
//		{
//		 "parameters": [
//		   {
//		     "key": "placeholder",
//		     "value": "value"
//		   }
//		 ],
//	  "secrets": [
//	    {
//	      "key": "placeholder",
//	      "value": "value"
//		  }
//		 ]
//		}
type Parameters struct {
	Metadata

	// Parameters is a list of key/value pairs.
	Parameters []Parameter `json:"parameters,omitempty"`

	// Secrets is a list of key/value pairs where value is write-only.
	// After being set, the API returns a secret reference instead of the actual value.
	// Setting a new value for the same key will invalidate the previous secret reference.
	// Secret values are redacted from audit and device logs.
	Secrets []Parameter `json:"secrets,omitempty"`
}

// Parameter defines a key/value pair.
type Parameter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
