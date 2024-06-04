// Copyright 2024 qbee.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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

// RunAgentManager is a manager for running the agent on devices.
type RunAgentManager struct {
	// client is the RunAgentManager client.
	client *Client
	// allowFailures allows failures.
	allowFailures bool
}

// NewRunAgentManager returns a new RunAgentManager.
func NewRunAgentManager() *RunAgentManager {
	return &RunAgentManager{}
}

// WithClient sets the client
func (m *RunAgentManager) WithClient(client *Client) *RunAgentManager {
	m.client = client
	return m
}

// WithAllowFailures sets the allowFailures flag
func (m *RunAgentManager) WithAllowFailures(allow bool) *RunAgentManager {
	m.allowFailures = allow
	return m
}

// RunAgentDevice runs the agent on the device.
func (m *RunAgentManager) RunAgentDevice(ctx context.Context, deviceID string) error {
	err := m.runAgent(ctx, deviceID)
	if err == nil {
		fmt.Printf("Successfully triggered agent on device %s\n", deviceID)
		return nil
	}
	if m.allowFailures {
		fmt.Printf("Failed to run agent on device %s: %s\n", deviceID, err.Error())
		return nil
	}
	return fmt.Errorf("failed to run agent on device %s: %w", deviceID, err)
}

// RunAgentGroup runs the agent on all devices in the group.
func (m *RunAgentManager) RunAgentGroup(ctx context.Context, groupID string) error {

	query := InventoryListQuery{
		Search: InventoryListSearch{
			Ancestors: []string{groupID},
		},
	}

	return m.runAgentMultiple(ctx, query)

}

// RunAgentTag runs the agent on all devices in the group with the specified tag.
func (m *RunAgentManager) RunAgentTag(ctx context.Context, tag string) error {
	query := InventoryListQuery{
		Search: InventoryListSearch{
			Tags: []string{tag},
		},
	}

	return m.runAgentMultiple(ctx, query)

}

// runAgentMultiple runs the agent on multiple devices.
func (m *RunAgentManager) runAgentMultiple(ctx context.Context, query InventoryListQuery) error {
	devices, err := m.client.ListDeviceInventory(ctx, query)

	if err != nil {
		return err
	}

	if len(devices.Items) == 0 {
		return fmt.Errorf("no devices found")
	}

	for _, device := range devices.Items {
		var err error
		if err = m.runAgent(ctx, device.NodeID); err == nil {
			fmt.Printf("Successfully triggered agent on device %s\n", device.NodeID)
			continue
		}
		if !m.allowFailures {
			return fmt.Errorf("failed to run agent on device %s: %w", device.NodeID, err)
		}
		fmt.Printf("Failed to run agent on device %s: %s\n", device.NodeID, err.Error())
	}
	return nil
}

// runAgent runs the agent on the device.
func (m *RunAgentManager) runAgent(ctx context.Context, deviceID string) error {
	path := fmt.Sprintf("/api/v2/device/%s/run-agent", deviceID)

	return m.client.Call(ctx, http.MethodGet, path, nil, nil)
}
