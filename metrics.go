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
)

// FilesystemMount contains the summary of a filesystem mount.
type FilesystemMount struct {
	// Path is the mount point of the filesystem.
	Path string `json:"path"`

	// Available is the amount of available (free) disk space in kilobytes.
	Available int64 `json:"avail"`

	// Utilization is the percentage of used disk space.
	Utilization int64 `json:"util"`
}

// FilesystemSummary contains the summary of the filesystems on a device.
type FilesystemSummary struct {
	// Reported contains unix timestamp of when the metric was recorded.
	Reported int64 `json:"reported"`

	// Mounts contains the summary of the filesystem mounts.
	Mounts []FilesystemMount `json:"mounts"`
}

// CPUSummary contains the summary of the CPU metrics on a device.
type CPUSummary struct {
	// Reported contains unix timestamp of when the metric was recorded.
	Reported int64 `json:"reported"`

	// User is the percentage of CPU time spent in the user space.
	User float64 `json:"user"`

	// System is the percentage of CPU time spent in the kernel space.
	System float64 `json:"system"`

	// IO is the percentage of CPU time spent on waiting for I/O.
	IO float64 `json:"io"`
}

// LoadSummary contains the summary of the load metrics on a device.
type LoadSummary struct {
	// Reported contains unix timestamp of when the metric was recorded.
	Reported int64 `json:"reported"`

	// Minute1 average system load over 1 minute.
	Minute1 float64 `json:"1min"`

	// Minute1 average system load over 5 minutes.
	Minute5 float64 `json:"5min"`

	// Minute1 average system load over 15 minutes.
	Minute15 float64 `json:"15min"`
}

// MemorySummary contains the summary of the memory metrics on a device.
type MemorySummary struct {
	// Reported contains unix timestamp of when the metric was recorded.
	Reported int64 `json:"reported"`

	// Available is the amount of available (free) memory in kilobytes.
	Available int64 `json:"avail"`

	// Used is the percentage of memory being used.
	Utilization int64 `json:"util"`

	// SwapUtilization is the percentage of swap being used.
	SwapUtilization int64 `json:"swap"`
}

// MetricsSummary contains the summary of the metrics on a device.
type MetricsSummary struct {
	Filesystem FilesystemSummary `json:"filesystem"`
	CPU        CPUSummary        `json:"cpu"`
	Load       LoadSummary       `json:"load"`
	Memory     MemorySummary     `json:"memory"`
}

// GetLatestMetrics returns the latest metrics of a device (not older than 1 hour).
func (cli *Client) GetLatestMetrics(ctx context.Context, deviceID string) (*MetricsSummary, error) {
	metricsSummary := new(MetricsSummary)

	path := fmt.Sprintf("/api/v2/device/%s/metrics/latest", deviceID)

	if err := cli.Call(ctx, http.MethodGet, path, nil, metricsSummary); err != nil {
		return nil, err
	}

	return metricsSummary, nil
}
