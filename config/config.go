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

// EntityType is used to distinguish between node and tag config.
type EntityType string

const (
	// EntityTypeNode represents a node entity.
	EntityTypeNode EntityType = "node"

	// EntityTypeTag represents a tag entity.
	EntityTypeTag EntityType = "tag"
)

// EntityConfigScope is used to distinguish between different scopes of config.
type EntityConfigScope string

const (
	// EntityConfigScopeAll returns final calculated config (incl. ancestors and tags)
	EntityConfigScopeAll EntityConfigScope = "all"

	// EntityConfigScopeOwn returns only config for the entity itself (no ancestors or tags)
	EntityConfigScopeOwn EntityConfigScope = "own"
)

// Config contains entity's configuration bundles
type Config struct {
	// EntityID is either a nodeID for EntityTypeNode or tag value for EntityTypeTag
	EntityID string `json:"id"`

	// Type defines entity type ID relevant to above EntityID
	Type EntityType `json:"type"`

	// CommitID of the most recent commit affecting the config's contents
	CommitID string `json:"commit_id"`

	// CommitCreated is a creation timestamp of the commit used to determine which Config has the most recent changes
	// in the chain of configs when we calculate an active config for an entity.
	// This is stored in nanosecond resolution.
	CommitCreated int64 `json:"commit_created"`

	// Bundles contains a list of strings representing configuration bundles
	Bundles BundleNames `json:"bundles"`

	// BundleData contain configuration data for bundles in the Bundles list
	BundleData BundleData `json:"bundle_data"`
}
