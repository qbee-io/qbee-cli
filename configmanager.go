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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"go.qbee.io/client/config"
)

// ConfigurationManager is a configuration for uploads to config entities
type ConfigurationManager struct {
	// client is the ConfiguratinManager client.
	client *Client

	// templateParameters are the template parameters.
	templateParameters map[string]string

	// entityType is the entity type.
	entityType config.EntityType

	// entityConfigScope is the entity scope.
	entityConfigScope config.EntityConfigScope
}

// NewConfigurationManager creates a new ConfigurationManager.
func NewConfigurationManager() *ConfigurationManager {
	return &ConfigurationManager{}
}

// WithClient sets the client.
func (c *ConfigurationManager) WithClient(client *Client) *ConfigurationManager {
	c.client = client
	return c
}

// WithTemplateParameters sets the template parameters.
func (c *ConfigurationManager) WithTemplateParameters(params map[string]string) *ConfigurationManager {
	c.templateParameters = params
	return c
}

// WithEntityType sets the entity type.
func (c *ConfigurationManager) WithEntityType(entityType config.EntityType) *ConfigurationManager {
	c.entityType = entityType
	return c
}

// WithEntityConfigScope sets the entity config scope.
func (c *ConfigurationManager) WithEntityConfigScope(entityConfigScope config.EntityConfigScope) *ConfigurationManager {
	c.entityConfigScope = entityConfigScope
	return c
}

// Save saves the configuration for the entity type.
func (c *ConfigurationManager) Save(ctx context.Context, target, bundleName, configFile string) error {

	var configBytes []byte
	var err error
	var configData any

	if configBytes, err = os.ReadFile(configFile); err != nil {
		return err
	}

	if c.templateParameters != nil {
		if configBytes, err = c.expandTemplate(configBytes); err != nil {
			return err
		}
	}

	if err = json.Unmarshal(configBytes, &configData); err != nil {
		return err
	}

	configChange := Change{
		BundleName: bundleName,
		Config:     configData,
	}

	if c.entityType == config.EntityTypeNode {
		configChange.NodeID = target
	} else if c.entityType == config.EntityTypeTag {
		configChange.Tag = target
	} else {
		return fmt.Errorf("invalid entity type: %s", c.entityType)
	}

	if _, err = c.client.CreateConfigurationChange(ctx, configChange); err != nil {
		return err
	}

	return nil

}

// Commit commits the configuration if there are uncommitted changes.
func (c *ConfigurationManager) Commit(ctx context.Context, commitMessage string) error {
	if _, err := c.client.CommitConfiguration(ctx, commitMessage); err != nil {
		return err
	}
	return nil
}

// GetConfig returns the active configuration for the entity.
func (c *ConfigurationManager) GetConfig(
	ctx context.Context,
	entityID string,
) (*config.Config, error) {

	return c.client.GetActiveConfig(ctx, c.entityType, entityID, c.entityConfigScope)
}

// expandTemplate expands the template using built in Golang templating.
func (c *ConfigurationManager) expandTemplate(configBytes []byte) ([]byte, error) {

	tmpl, err := template.New("config-template").Parse(string(configBytes))

	// Produce error on missing keys.
	tmpl.Option("missingkey=error")

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	if err := tmpl.Execute(buf, c.templateParameters); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
