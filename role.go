// Copyright 2025 qbee.io
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

	// CreatedAt is the timestamp of the creation of the role.
	CreatedAt int64 `json:"created_at"`

	// CreatedBy is the user information of the user that created the role.
	CreatedBy *UserBaseInfo `json:"created_by"`

	// UpdatedAt is the timestamp of the last update of the role.
	UpdatedAt int64 `json:"updated_at,omitempty"`

	// UpdatedBy is the user information of the user that last updated the role
	UpdatedBy *UserBaseInfo `json:"updated_by,omitempty"`

	// list of users using this role
	UsedBy []UserBaseInfo `json:"used_by"`
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
