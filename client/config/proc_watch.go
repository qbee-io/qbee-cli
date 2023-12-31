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

// ProcessWatchBundle defines name for the process watch bundle.
const ProcessWatchBundle Bundle = "proc_watch"

// ProcessWatch ensures running process are running (or not).
//
// Example payload:
//
//	{
//	  "processes": [
//	   {
//	     "name": "presentProcess",
//	     "policy": "Present",
//	     "command": "start.sh"
//	   },
//	   {
//	     "name": "absentProcess",
//	     "policy": "Absent",
//	     "command": "stop.sh"
//	   }
//	 ]
//	}
type ProcessWatch struct {
	Metadata

	// Processes to watch.
	Processes []ProcessWatcher `json:"processes,omitempty"`
}

// ProcessWatcher defines a watcher for a process.
type ProcessWatcher struct {
	// Name of the process to watch.
	Name string `json:"name,omitempty"`

	// Policy for the process.
	Policy ProcessPolicy `json:"policy,omitempty"`

	// Command to use to get the process in the expected state.
	// For:
	// - ProcessPresent it should be a start command,
	// - ProcessAbsent it should be a stop command.
	Command string `json:"command,omitempty"`
}

// ProcessPolicy defines a policy for a process.
type ProcessPolicy string

const (
	// ProcessPresent ensures that the process is running.
	ProcessPresent ProcessPolicy = "Present"

	// ProcessAbsent ensures that the process is not running.
	ProcessAbsent ProcessPolicy = "Absent"
)
