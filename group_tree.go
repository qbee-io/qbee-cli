package client

import (
	"context"
	"net/http"
)

// GroupTreeAction represents an action on the group tree.
type GroupTreeAction string

// GroupTreeAction types enumeration.
const (
	TreeActionCreate GroupTreeAction = "create"
	TreeActionRename GroupTreeAction = "rename"
	TreeActionUpdate GroupTreeAction = "update"
	TreeActionMove   GroupTreeAction = "move"
	TreeActionDelete GroupTreeAction = "delete"
)

// System node IDs
const (
	// RootNodeID is the ID of the root node.
	// Root node cannot be moved or deleted.
	RootNodeID = "root"

	// UnassignedGroupNodeID is the ID of the unassigned group node.
	// Unassigned group node cannot be moved or deleted.
	UnassignedGroupNodeID = "unassigned_group"
)

// GroupTreeChangeData represents the data for a group tree change.
type GroupTreeChangeData struct {
	ParentID string `json:"parent_id"`
	NodeID   string `json:"node_id"`

	// Type is the type of the group (device or group).
	// This parameter is used only for moving of existing groups.
	Type NodeType `json:"type,omitempty"`

	// OldParentID is the ID of the parent group before the move.
	// This parameter is used only for moving groups.
	OldParentID string `json:"old_parent_id,omitempty"`

	// Title is the title of the new group.
	// This parameter is used only for creating new groups.
	Title string `json:"title,omitempty"`

	// Tags is the list of tags to set on the node.
	// This parameter is used only for updates.
	Tags []string `json:"tags,omitempty"`
}

// GroupTree represents the device group tree.
type GroupTree struct {
	TotalDevices int           `json:"total_devices"`
	Tree         GroupTreeNode `json:"tree"`
}

// GroupTreeNode represents a node in the device group tree.
type GroupTreeNode struct {
	ID     string          `json:"id"`
	NodeID string          `json:"node_id"`
	Title  string          `json:"title"`
	Type   NodeType        `json:"type"`
	Tags   []string        `json:"tags,omitempty"`
	Nodes  []GroupTreeNode `json:"nodes"`

	// device only attributes
	PublicKeyDigest string `json:"pub_key_digest,omitempty"`
	DeviceCommitSHA string `json:"device_commit_sha,omitempty"`

	Status           string           `json:"status,omitempty"`
	ConfigPropagated bool             `json:"config_propagated,omitempty"`
	AgentInterval    int              `json:"agentinterval,omitempty"`
	LastReported     int64            `json:"last_reported,omitempty"`
	Attributes       DeviceAttributes `json:"attributes"`

	// internal attributes
	ParentID string `json:"-"`
	Hostname string `json:"-"`
}

// GroupTreeChange represents a change in the device tree.
// Data is specific to the action:
// - create: parent_id, node_id, title, type
// - rename: node_id, title, type=group (only groups can be renamed)
// - move: parent_id, node_id, old_parent_id, type=device (only devices can be moved)
// - delete: node_id, type=group (only empty groups can be deleted)
type GroupTreeChange struct {
	Action GroupTreeAction     `json:"action"`
	Data   GroupTreeChangeData `json:"data"`
}

const grouptreePath = "/api/v2/grouptree"

// GroupTreeRequest is the request to update the group tree.
type GroupTreeRequest struct {
	Changes []GroupTreeChange `json:"changes"`
}

// GroupTreeUpdate updates the group tree.
func (cli *Client) GroupTreeUpdate(ctx context.Context, req GroupTreeRequest) error {
	return cli.Call(ctx, http.MethodPut, grouptreePath, req, nil)
}

// GroupTreeGet returns the device group tree.
func (cli *Client) GroupTreeGet(ctx context.Context, skipUnassigned bool) (*GroupTree, error) {
	groupTree := new(GroupTree)

	path := grouptreePath
	if skipUnassigned {
		path += "?skip_unassigned=true"
	}

	if err := cli.Call(ctx, http.MethodGet, path, nil, groupTree); err != nil {
		return nil, err
	}

	return groupTree, nil
}

// GroupTreeGetNode returns the node from device group tree.
func (cli *Client) GroupTreeGetNode(ctx context.Context, nodeID string) (*NodeInfo, error) {
	nodeInfo := new(NodeInfo)

	path := grouptreePath + "/" + nodeID

	if err := cli.Call(ctx, http.MethodGet, path, nil, nodeInfo); err != nil {
		return nil, err
	}

	return nodeInfo, nil
}

// GroupTreeSetTagsRequest is the request to set tags for a node.
type GroupTreeSetTagsRequest struct {
	Tags []string `json:"tags"`
}

// GroupTreeSetTags sets tags for a node.
func (cli *Client) GroupTreeSetTags(ctx context.Context, nodeID string, tags []string) error {
	path := grouptreePath + "/" + nodeID

	req := GroupTreeSetTagsRequest{
		Tags: tags,
	}

	return cli.Call(ctx, http.MethodPatch, path, req, nil)
}
