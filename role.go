package client

import (
	"context"
	"fmt"
	"net/http"
)

// Permission represents a permission in the system.
type Permission string

// Role represents a role in the system.
type Role struct {
	// ID is the unique identifier of the role.
	ID string `json:"id"`

	// Name is the short name of the role.
	Name string `json:"name"`

	// Description is the optional description of the role.
	Description string `json:"description"`

	// Policies is the list of policies that are assigned to this role.
	Policies []RolePolicy `json:"policies"`
}

// RolePolicy represents a policy that can be assigned to a role.
type RolePolicy struct {
	// Permission is the permission that is granted by this policy.
	Permission Permission `json:"permission"`

	// Resources is the list of resources that are affected by this policy.
	Resources []string `json:"resources,omitempty"`
}

const rolePath = "/api/v2/role"

// CreateRole creates a new role.
func (cli *Client) CreateRole(ctx context.Context, role Role) (*Role, error) {
	createdRole := new(Role)

	if err := cli.Call(ctx, http.MethodPost, rolePath, role, createdRole); err != nil {
		return nil, err
	}

	return createdRole, nil
}

// UpdateRole updates an existing role.
func (cli *Client) UpdateRole(ctx context.Context, role Role) (*Role, error) {
	path := fmt.Sprintf("%s/%s", rolePath, role.ID)

	updatedRole := new(Role)

	if err := cli.Call(ctx, http.MethodPut, path, role, updatedRole); err != nil {
		return nil, err
	}

	return updatedRole, nil
}

// DeleteRole deletes a role by its ID.
func (cli *Client) DeleteRole(ctx context.Context, roleID string) error {
	path := fmt.Sprintf("%s/%s", rolePath, roleID)

	if err := cli.Call(ctx, http.MethodDelete, path, nil, nil); err != nil {
		return err
	}

	return nil
}

// GetRole returns a role by its ID.
func (cli *Client) GetRole(ctx context.Context, roleID string) (*Role, error) {
	path := fmt.Sprintf("%s/%s", rolePath, roleID)

	role := new(Role)

	if err := cli.Call(ctx, http.MethodGet, path, nil, role); err != nil {
		return nil, err
	}

	return role, nil
}

const roleListPath = "/api/v2/roleslist"

// ListRoles returns a list of all roles in the system.
func (cli *Client) ListRoles(ctx context.Context) ([]Role, error) {
	var response struct {
		Status string `json:"status"`
		Total  int    `json:"total"`
		Roles  []Role `json:"items"`
	}

	if err := cli.Call(ctx, http.MethodGet, roleListPath, nil, &response); err != nil {
		return nil, err
	}

	return response.Roles, nil
}
