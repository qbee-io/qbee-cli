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
