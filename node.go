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
