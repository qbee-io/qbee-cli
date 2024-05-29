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
	"fmt"
	"net/http"

	"go.qbee.io/client/config"
)

// ConfigResponse is a response for the configuration request.
type ConfigResponse struct {
	// Status is the status of the request.
	Status string `json:"status"`

	// Config is the configuration of the requested entity.
	Config *config.Config `json:"config"`
}

// GetActiveConfig returns the active configuration for the entity.
func (cli *Client) GetActiveConfig(
	ctx context.Context,
	entityType config.EntityType,
	entityID string,
	scope config.EntityConfigScope,
) (*config.Config, error) {
	path := fmt.Sprintf("/api/v2/config/%s/%s", entityType, entityID)

	if scope != "" {
		path += fmt.Sprintf("?scope=%s", scope)
	}

	response := new(ConfigResponse)

	if err := cli.Call(ctx, http.MethodGet, path, nil, response); err != nil {
		return nil, err
	}

	return response.Config, nil
}

// GetConfigPreview returns the configuration preview (with uncommitted changes) for the entity.
func (cli *Client) GetConfigPreview(
	ctx context.Context,
	entityType config.EntityType,
	entityID string,
	scope config.EntityConfigScope,
) (*config.Config, error) {
	path := fmt.Sprintf("/api/v2/configpreview/%s/%s", entityType, entityID)

	if scope != "" {
		path += fmt.Sprintf("?scope=%s", scope)
	}

	response := new(ConfigResponse)

	if err := cli.Call(ctx, http.MethodGet, path, nil, response); err != nil {
		return nil, err
	}

	return response.Config, nil
}

// UploadConfig uploads the configuration for the entity.
type ConfigPayload struct {
	NodeID   string `json:"node_id,omitempty"`
	FormType string `json:"formtype,omitempty"`
	Config   any    `json:"config,omitempty"`
	Tag      string `json:"tag,omitempty"`
}

// configPath is the path for the config request.
const configPath = "/api/v2/change"

// UploadConfig uploads the configuration for the entity.
func (cli *Client) UploadConfig(
	ctx context.Context,
	payload *ConfigPayload,
) error {

	if err := cli.Call(ctx, http.MethodPost, configPath, payload, nil); err != nil {
		return err
	}

	return nil
}

// CommitPayload is a payload for the commit request.
type CommitPayload struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

// commitPath is the path for the commit request.
const commitPath = "/api/v2/commit"

// CommitConfig commits the configuration for the entity.
func (cli *Client) CommitConfig(
	ctx context.Context,
	commitMessage string,
) error {

	payload := &CommitPayload{
		Action:  "commit",
		Message: commitMessage,
	}

	if err := cli.Call(ctx, http.MethodPost, commitPath, payload, nil); err != nil {
		return err
	}

	return nil
}
