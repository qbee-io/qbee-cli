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

// Container defines a container instance.
type Container struct {
	// Name used by the container.
	Name string `json:"name"`

	// Image used by the container.
	Image string `json:"image"`

	// PodmanArgs defines command line arguments for "podman run".
	DockerArgs string `json:"docker_args"`

	// EnvFile defines an env file (from file manager) to be used inside container.
	EnvFile string `json:"env_file"`

	// Command to be executed in the container.
	Command string `json:"command"`

	// PreCondition defines a condition that must be met before the container is started.
	PreCondition string `json:"pre_condition"`
}

// RegistryAuth defines credentials for podman registry authentication.
type RegistryAuth struct {
	// Server hostname of the registry.
	Server string `json:"server"`

	// Username for the registry.
	Username string `json:"username"`

	// Password for the Username.
	Password string `json:"password"`
}
