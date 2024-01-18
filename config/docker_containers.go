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

// DockerContainersBundle defines name for the docker containers bundle.
const DockerContainersBundle Bundle = "docker_containers"

// DockerContainers controls docker containers running in the system.
//
// Example payload:
//
//	{
//		"items": [
//		  {
//	     "name": "container-a",
//	     "image": "debian:stable",
//	     "docker_args": "-v /path/to/data-volume:/data --hostname my-hostname",
//	     "env_file": "/my-directory/my-envfile",
//	     "command": "echo 'hello world!'"
//		  }
//		],
//	 "registry_auths": [
//	   {
//	      "server": "gcr.io",
//	      "username": "user",
//	      "password": "seCre7"
//	   }
//	 ]
//	}
type DockerContainers struct {
	Metadata

	// Containers to be running in the system.
	Containers []DockerContainer `json:"items,omitempty"`

	// RegistryAuths contains credentials to private docker registries.
	RegistryAuths []RegistryAuth `json:"registry_auths,omitempty"`
}

// DockerContainer defines a docker container instance.
type DockerContainer struct {
	// Name used by the container.
	Name string `json:"name"`

	// Image used by the container.
	Image string `json:"image"`

	// DockerArgs defines command line arguments for "docker run".
	DockerArgs string `json:"docker_args"`

	// EnvFile defines an env file (from file manager) to be used inside container.
	EnvFile string `json:"env_file"`

	// Command to be executed in the container.
	Command string `json:"command"`
}

// RegistryAuth defines credentials for docker registry authentication.
type RegistryAuth struct {
	// Server hostname of the registry.
	Server string `json:"server"`

	// Username for the registry.
	Username string `json:"username"`

	// Password for the Username.
	Password string `json:"password"`
}
