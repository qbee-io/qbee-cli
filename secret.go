package client

import (
	"context"
	"net/http"
)

const secretPath = "/api/v2/secret"

// CreateSecretRequest is the request body for creating a secret.
type CreateSecretRequest struct {
	// Value is the plain text value of the secret.
	Value string `json:"value"`
}

// CreateSecretResponse is the response from the server when creating a secret.
type CreateSecretResponse struct {
	// SecretRef is the reference to the created secret.
	SecretRef string `json:"secret_ref"`
}

// CreateSecret stores the provided plain text value as a secret and returns its reference.
func (cli *Client) CreateSecret(ctx context.Context, value string) (string, error) {
	request := CreateSecretRequest{Value: value}
	response := new(CreateSecretResponse)

	if err := cli.Call(ctx, http.MethodPost, secretPath, request, response); err != nil {
		return "", err
	}

	return response.SecretRef, nil
}
