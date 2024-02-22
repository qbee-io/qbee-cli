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

package client

import "go.qbee.io/client/config"

// NodeType is the type of node.
type NodeType string

// Available node types.
const (
	NodeTypeDevice NodeType = "device"
	NodeTypeGroup  NodeType = "group"
)

// NodeInfo contains information about a node.
// This is used for the device tree API.
type NodeInfo struct {
	NodeID          string     `json:"node_id"`
	PublicKeyDigest string     `json:"pub_key_digest"`
	Type            NodeType   `json:"type"`
	Ancestors       []string   `json:"ancestors"`
	Title           string     `json:"title"`
	Tags            []string   `json:"tags,omitempty"`
	Nodes           []NodeInfo `json:"nodes,omitempty"`

	// Device specific fields
	UUID             string            `json:"uuid,omitempty"`
	Status           string            `json:"status,omitempty"`
	DeviceCommitSHA  string            `json:"device_commit_sha,omitempty"`
	Attributes       *DeviceAttributes `json:"attributes,omitempty"`
	ConfigPropagated bool              `json:"config_propagated,omitempty"`
	AgentInterval    int               `json:"agentinterval,omitempty"`
	LastReported     int64             `json:"last_reported,omitempty"`
	System           *SystemInfo       `json:"system,omitempty"`
	Settings         *config.Settings  `json:"settings,omitempty"`
	PushedConfig     *config.Pushed    `json:"pushed_config,omitempty"`
}
