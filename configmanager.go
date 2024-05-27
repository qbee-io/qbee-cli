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

import "context"

type ConfiguratinManager struct {
	// client is the ConfiguratinManager client.
	client *Client
}

// uploadConfiguration uploads configuration of any type to the entity.
func (c *ConfiguratinManager) uploadConfiguration(ctx context.Context, entityID string, config []byte) error {

	return nil
}

// commitConfiguration commits the configuration for the entity.
func (c *ConfiguratinManager) commitConfiguration(ctx context.Context, entityID string) error {

	return nil
}

// UploadAndCommitConfiguration uploads and commits the configuration for the entity.

func (c *ConfiguratinManager) UploadAndCommitConfiguration(ctx context.Context, commitMessage, entityID string, config []byte) error {

	if err := c.uploadConfiguration(ctx, entityID, config); err != nil {
		return err
	}

	if err := c.commitConfiguration(ctx, commitMessage); err != nil {
		return err
	}

	return nil
}
