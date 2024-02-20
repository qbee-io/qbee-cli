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
