package client

import (
	"context"
	"net/http"
)

// BootstrapKey represents a bootstrap key in the system.
type BootstrapKey struct {
	// ID is the actual bootstrap key.
	ID string `json:"id"`
	// GroupID is the group ID of associated with the bootstrap key.
	GroupID string `json:"group_id,omitempty"`
	// AutoAccept indicates whether the bootstrap key is auto accepted.
	AutoAccept bool `json:"auto_accept"`
}

// BootstrapKeyResponse represents the response from the server when creating or updating a bootstrap key.
type BootstrapKeyResponse map[string]BootstrapKey

const bootstrapKeyPath = "/api/v2/bootstrapkey"

// First returns the first bootstrap key in the response.
func (keys BootstrapKeyResponse) First() *BootstrapKey {
	for keyID, key := range keys {
		key.ID = keyID
		return &key
	}

	return nil
}

// NewBootstrapKey returns a new bootstrap key.
func (cli *Client) NewBootstrapKey(ctx context.Context) (*BootstrapKey, error) {
	response := make(BootstrapKeyResponse)

	err := cli.Call(ctx, http.MethodPost, bootstrapKeyPath, nil, &response)
	if err != nil {
		return nil, err
	}

	return response.First(), nil
}

const bootstrapKeyListPath = "/api/v2/bootstrapkeylist"

// ListBootstrapKeys returns a list of all bootstrap keys in the system.
func (cli *Client) ListBootstrapKeys(ctx context.Context) (BootstrapKeyResponse, error) {
	response := make(BootstrapKeyResponse)

	err := cli.Call(ctx, http.MethodGet, bootstrapKeyListPath, nil, &response)
	if err != nil {
		return nil, err
	}

	// populate key IDs
	for keyID, key := range response {
		key.ID = keyID
		response[keyID] = key
	}

	return response, nil
}

// UpdateBootstrapKey in the system.
func (cli *Client) UpdateBoostrapKey(ctx context.Context, key *BootstrapKey) error {
	keyPath := bootstrapKeyPath + "/" + key.ID
	return cli.Call(ctx, http.MethodPut, keyPath, key, nil)
}

// DeleteBootstrapKey from the system.
func (cli *Client) DeleteBootstrapKey(ctx context.Context, keyID string) error {
	keyPath := bootstrapKeyPath + "/" + keyID

	return cli.Call(ctx, http.MethodDelete, keyPath, nil, nil)
}

// GetBootstrapKey returns a bootstrap key by its ID.
func (cli *Client) GetBootstrapKey(ctx context.Context, keyID string) (*BootstrapKey, error) {
	keyPath := bootstrapKeyPath + "/" + keyID
	response := make(BootstrapKeyResponse)

	err := cli.Call(ctx, http.MethodGet, keyPath, nil, &response)
	if err != nil {
		return nil, err
	}

	return response.First(), nil
}
