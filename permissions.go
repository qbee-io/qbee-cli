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

const (
	// Analysis Section
	AnalysisRead Permission = "analysis:read"

	// Audit Section
	AuditRead Permission = "audit:read"

	// Billing Section
	BillingManage Permission = "billing:manage"
	BillingRead   Permission = "billing:read"

	// Bootstrap Keys Section
	BootstrapKeysManage Permission = "bootstrap-keys:manage"
	BootstrapKeysRead   Permission = "bootstrap-keys:read"

	// Company Section
	CompanyManage Permission = "company:manage"
	CompanyRead   Permission = "company:read"

	// Configuration Section
	ConfigurationManage Permission = "configuration:manage"
	ConfigurationRead   Permission = "configuration:read"

	// CVE Section
	CVEManage Permission = "cve:manage"
	CVERead   Permission = "cve:read"

	// Devices Section
	DeviceApprove Permission = "device:approve"
	DeviceManage  Permission = "device:manage"
	DeviceRead    Permission = "device:read"

	// Files Section
	FilesManage Permission = "files:manage"
	FilesRead   Permission = "files:read"

	// Remote Access Section
	RemoteAccessConnect Permission = "remote-access:connect"
	RemoteAccessManage  Permission = "remote-access:manage"

	// Reports Section
	ReportsAcknowledge Permission = "reports:acknowledge"
	ReportsRead        Permission = "reports:read"

	// Roles Section
	RolesManage Permission = "roles:manage"
	RolesRead   Permission = "roles:read"

	// Users Section
	UsersManage Permission = "users:manage"
	UsersRead   Permission = "users:read"
)
