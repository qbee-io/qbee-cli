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

const DockerContainersBundle Bundle = "docker_containers"

// DockerContainers controls docker containers running in the system.
//
// Example payload:
//
//	{
//		"items": [
//		  {
//			"name": "container-a",
//			"image": "debian:stable",
//			"podman_args": "-v /path/to/data-volume:/data --hostname my-hostname",
//			"env_file": "/my-directory/my-envfile",
//			"command": "echo hello world!"
//		  }
//		],
//		"registry_auths": [
//		  {
//			"server": "gcr.io",
//			"username": "user",
//			"password": "seCre7"
//		  }
//		]
//	}
type DockerContainers struct {
	Metadata `bson:"-,inline"`

	// Containers to be running in the system.
	Containers []Container `json:"items,omitempty" bson:"items,omitempty"`

	// RegistryAuths contains credentials to private docker registries.
	RegistryAuths []RegistryAuth `json:"registry_auths,omitempty" bson:"registry_auths,omitempty"`
}
