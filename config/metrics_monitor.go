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

package config

// MetricsMonitorBundle defines name for the metrics_monitor bundle.
const MetricsMonitorBundle Bundle = "metrics_monitor"

// MetricsMonitor configures on-agent metrics monitoring.
//
// Example payload:
//
//	{
//	  "metrics": [
//	 	{
//	 		"value": "cpu:user",
//			"threshold": 20.0
//	 	},
//		{
//	 		"value": "filesystem:use",
//			"threshold": 60.0,
//			"id": "/data"
//	 	},
//
//	  ]
//	}
type MetricsMonitor struct {
	Metadata

	Metrics []MetricMonitor `json:"metrics"`
}

// MetricMonitor defines monitor for a single metric.
type MetricMonitor struct {
	// Value of the metric (enum defined in the JSON schema).
	Value string `json:"value"`

	// Threshold above which a warning will be created by the device.
	Threshold float64 `json:"threshold"`

	// ID of the resource (e.g. filesystem mount point)
	ID string `json:"id,omitempty"`
}
