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

// SSHKeysBundle defines name for the SSH keys management bundle.
const SSHKeysBundle Bundle = "sshkeys"

// SSHKeys adds or removes authorized SSH keys for users.
//
// Example payload:
//
//	{
//	 "users": [
//	   {
//	     "username": "test",
//	     "userkeys": [
//	       "key1",
//	       "key2"
//	     ]
//	   }
//	 ]
//	}
type SSHKeys struct {
	Metadata

	// Users to add or remove SSH keys for.
	Users []SSHKey `json:"users,omitempty"`
}

// SSHKey defines an SSH key to be added to a user.
type SSHKey struct {
	// Username of the user for which the SSH key is added.
	Username string `json:"username"`

	// UserKeys are SSH keys to be added to the user.
	Keys []string `json:"userkeys"`
}
