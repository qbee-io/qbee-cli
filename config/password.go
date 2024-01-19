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

// PasswordBundle defines name for the passwords bundle.
const PasswordBundle Bundle = "password"

// Password bundle sets passwords for existing users.
//
// Example payload:
//
//	{
//	 "users": [
//	   {
//	     "username": "piotr",
//	     "passwordhash": "$6$EMNbdq1ZkOAZSpFt$t6Ei4J11Ybip1A51sbBPTtQEVcFPPPUs.Q9nle4FenvrId4fLr8douwE3lbgWZGK.LIPeVmmFrTxYJ0QoYkFT."
//	   }
//	 ]
//	}
type Password struct {
	Metadata

	Users []UserPassword `json:"users,omitempty"`
}

// UserPassword defines a user and password hash.
type UserPassword struct {
	// Username of the user for which the password hash is set.
	Username string `json:"username"`

	// PasswordHash is a password hash for the user.
	// See https://qbee.io/docs/qbee-password.html for more information.
	PasswordHash string `json:"passwordhash"`
}
