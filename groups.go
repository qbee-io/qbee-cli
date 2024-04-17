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

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GroupNode represents a group node in the device tree.
type GroupNode struct {
	// NodeID of the group.
	NodeID string `json:"node_id"`

	// ParentID is the node ID of the parent group.
	ParentID string `json:"parent_id"`

	// Tags assigned to the group.
	Tags []string `json:"tags"`

	// NodeType is the type of node (always "group").
	Type NodeType `json:"type"`

	// Title of the group.
	Title string `json:"title"`

	// Ancestors is an array of node IDs of the ancestors of the group.
	Ancestors []string `json:"ancestors"`

	// Updated is the unix timestamp of the last update.
	Updated int64 `json:"updated"`
}

// GroupListSearch defines search parameters for GroupListQuery.
type GroupListSearch struct {
	// ParentID - node ID of the parent group.
	ParentID string `json:"parent_id,omitempty"`

	// Tags - array of tags which MUST be preset (exact match)
	Tags []string `json:"tags,omitempty"`
}

// GroupListQuery defines query parameters for ChangeList.
type GroupListQuery struct {
	Search GroupListSearch
}

// String returns string representation of GroupListQuery which can be used as query string.
func (q GroupListQuery) String() (string, error) {
	values := make(url.Values)

	searchQueryJSON, err := json.Marshal(q.Search)
	if err != nil {
		return "", err
	}

	values.Set("search", string(searchQueryJSON))

	return values.Encode(), nil
}

// GroupsList returns a list of all group nodes matching provided query.
func (cli *Client) GroupsList(ctx context.Context, query *GroupListQuery) ([]GroupNode, error) {
	endpointURL := "/api/v2/groups"

	if query != nil {
		queryString, err := query.String()
		if err != nil {
			return nil, err
		}

		endpointURL = fmt.Sprintf("%s?%s", endpointURL, queryString)
	}

	nodes := make([]GroupNode, 0)

	if err := cli.Call(ctx, http.MethodGet, endpointURL, nil, &nodes); err != nil {
		return nil, err
	}

	return nodes, nil
}
