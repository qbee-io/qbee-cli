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
	"path"
	"time"

	"go.qbee.io/client/config"
	"go.qbee.io/client/types"
)

// ChangeStatus defines the status of a change.
type ChangeStatus string

const (
	// ChangeStatusNew is the status of a new change.
	ChangeStatusNew ChangeStatus = "new"

	// ChangeStatusCommitted is the status of a committed change.
	ChangeStatusCommitted ChangeStatus = "committed"
)

// Change in the system state.
type Change struct {
	// ID is the unique identifier of the change.
	ID string `json:"id,omitempty"`

	// Content is the type specific content of the change.
	Content any `json:"content"`

	// SHA is the pseudo-digest of the change.
	SHA string `json:"sha"`

	// Status is the current status of the change.
	Status ChangeStatus `json:"status"`

	// Type is the type of the change.
	Type string `json:"type"`

	// Created is the timestamp when the change was created.
	Created types.Timestamp `json:"created"`

	// Updated is the timestamp when the change was last updated.
	Updated types.Timestamp `json:"updated"`

	// UserID is the ID of the user who created the change.
	UserID string `json:"user_id"`

	// UserName is the first name and last name of the user who created the change.
	// This field is populated by the API handlers.
	UserName string `json:"user_name,omitempty"`

	// User contains the user who created the change.
	// This field is populated by the API handlers.
	User *UserBaseInfo `json:"user,omitempty"`

	// NodeInfo contains base information about the node that the change is associated with (if any).
	// This field is populated by the API handlers.
	NodeInfo *NodeInfo `json:"node,omitempty"`
}

type ChangeRequest struct {
	// NodeID is the ID of the node the change is for.
	NodeID string `json:"node_id,omitempty"`

	// Tag is the tag the change is for.
	Tag string `json:"tag,omitempty"`

	// BundleName is the name of the configuration bundle.
	BundleName config.Bundle `json:"formtype"`

	// Extend defines if the change is an extension of the existing configuration.
	Extend bool `json:"extend"`

	// Content is the configuration of the change.
	Content any `json:"config"`
}

const changePath = "/api/v2/change"

// CreateConfigurationChange in the system.
func (cli *Client) CreateConfigurationChange(ctx context.Context, change ChangeRequest) (*Change, error) {
	createdChange := new(Change)

	err := cli.Call(ctx, http.MethodPost, changePath, change, createdChange)
	if err != nil {
		return nil, err
	}

	return createdChange, nil
}

// GetConfigurationChange returns a single change by its SHA.
func (cli *Client) GetConfigurationChange(ctx context.Context, sha string) (*Change, error) {
	apiPath := path.Join(changePath, sha)

	change := new(Change)

	if err := cli.Call(ctx, http.MethodGet, apiPath, nil, change); err != nil {
		return nil, err
	}

	return change, nil
}

// DeleteConfigurationChange deletes a single change by its SHA.
func (cli *Client) DeleteConfigurationChange(ctx context.Context, sha string) error {
	apiPath := path.Join(changePath, sha)

	return cli.Call(ctx, http.MethodDelete, apiPath, nil, nil)
}

// ChangeListSearch defines search parameters for ChangeListQuery.
type ChangeListSearch struct {
	// UserID - user ID to search for (exact match).
	UserID string `json:"user_id,omitempty"`

	// GroupID - group ID to search for (exact match).
	GroupID string `json:"group_id,omitempty"`

	// Status - status to search for (exact match).
	Status ChangeStatus `json:"status,omitempty"`

	// FormType - form type to search for (exact match).
	FormType string `json:"form_type,omitempty"`

	// Type - type to search for (exact match).
	Type string `json:"type,omitempty"`
}

const listQueryDateTimeFormat = "2006-01-02 15:04:05"

// ChangeListQuery defines query parameters for ChangeList.
type ChangeListQuery struct {
	Search ChangeListSearch

	// CreatedStart defines start of the time range to search in, inclusive.
	// Start date UTC, format: YYYY-MM-DD hh:ii:ss. Default -1 month from now
	CreatedStart time.Time

	// CreatedEnd defines end of the time range to search in, exclusive.
	// End date UTC, format: YYYY-MM-DD hh:ii:ss. Default: now
	CreatedEnd time.Time

	// SortField defines field used to sort, 'created' by default.
	// Supported sort fields:
	// - id
	// - created
	Sort string

	// SortDirection defines sort direction, 'desc' by default.
	SortDirection string

	// Limit defines maximum number of records in result, default 3000, max 10000
	Limit int

	// Offset defines offset of the first record in result, default 0
	Offset int
}

// String returns string representation of CommitListQuery which can be used as query string.
func (q ChangeListQuery) String() (string, error) {
	values := make(url.Values)

	searchQueryJSON, err := json.Marshal(q.Search)
	if err != nil {
		return "", err
	}

	values.Set("search", string(searchQueryJSON))

	if !q.CreatedStart.IsZero() {
		values.Set("created_start", q.CreatedStart.Format(listQueryDateTimeFormat))
	}

	if !q.CreatedEnd.IsZero() {
		values.Set("created_end", q.CreatedEnd.Format(listQueryDateTimeFormat))
	}

	if q.Sort != "" {
		values.Set("sort", q.Sort)
	}

	if q.SortDirection != "" {
		values.Set("sort_direction", q.SortDirection)
	}

	if q.Limit != 0 {
		values.Set("limit", fmt.Sprintf("%d", q.Limit))
	}

	if q.Offset != 0 {
		values.Set("offset", fmt.Sprintf("%d", q.Offset))
	}

	return values.Encode(), nil
}

// ListUncommittedConfigurationChanges returns a list of uncommitted configuration changes.
func (cli *Client) ListUncommittedConfigurationChanges(ctx context.Context, query ChangeListQuery) ([]Change, error) {
	queryString, err := query.String()
	if err != nil {
		return nil, err
	}

	apiPath := "/api/v2/changelist?" + queryString

	var changes []Change

	if err = cli.Call(ctx, http.MethodGet, apiPath, nil, &changes); err != nil {
		return nil, err
	}

	return changes, nil
}

// DeleteUncommittedConfigurationChanges deletes all uncommitted configuration changes.
func (cli *Client) DeleteUncommittedConfigurationChanges(ctx context.Context) error {
	const apiPath = "/api/v2/changes"

	return cli.Call(ctx, http.MethodDelete, apiPath, nil, nil)
}
