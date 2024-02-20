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

package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"go.qbee.io/client"
)

var devicesCommand = Command{
	Description: "Device management commands",
	SubCommands: map[string]Command{
		"list": devicesListSubcommand,
	},
}

const (
	devicesListQueryOption  = "query"
	devicesListLimitOption  = "limit"
	devicesListPageOption   = "page"
	devicesListAsJSONOption = "json"
)

var devicesListSubcommand = Command{
	Description: "List devices",
	Options: []Option{
		{
			Name:  devicesListQueryOption,
			Short: "q",
			Help:  "Query devices by name",
		},
		{
			Name:    devicesListLimitOption,
			Short:   "n",
			Help:    "Items per page",
			Default: "10",
		},
		{
			Name:    devicesListPageOption,
			Short:   "p",
			Help:    "Results page",
			Default: "0",
		},
		{
			Name:  devicesListAsJSONOption,
			Short: "j",
			Help:  "Output as JSON",
			Flag:  "true",
		},
	},
	Target: func(opts Options) error {
		ctx := context.Background()

		cli, err := client.LoginGetAuthenticatedClient(ctx)
		if err != nil {
			return err
		}

		query := client.InventoryListQuery{
			Search: client.InventoryListSearch{
				Title: opts[devicesListQueryOption],
			},
			SortField:     "title",
			SortDirection: client.SortDirectionAsc,
			ReportType:    "short",
		}

		if query.ItemsPerPage, err = strconv.Atoi(opts[devicesListLimitOption]); err != nil {
			return fmt.Errorf("invalid limit: %s", opts[devicesListLimitOption])
		}

		if query.Offset, err = strconv.Atoi(opts[devicesListPageOption]); err != nil {
			return fmt.Errorf("invalid offset: %s", opts[devicesListPageOption])
		} else {
			query.Offset *= query.ItemsPerPage
		}

		var devices *client.InventoryListResponse
		if devices, err = cli.ListDeviceInventory(ctx, query); err != nil {
			return err
		}

		if opts[devicesListAsJSONOption] == "true" {
			return json.NewEncoder(os.Stdout).Encode(devices)
		}

		tabWriter := tabwriter.NewWriter(os.Stdout, 14, 1, 1, ' ', 0)
		csvWriter := csv.NewWriter(tabWriter)
		csvWriter.Comma = '\t'
		_ = csvWriter.Write([]string{"Device ID", "Status", "Name"})
		for _, device := range devices.Items {
			_ = csvWriter.Write([]string{device.NodeID, device.Status, device.Title})
		}
		csvWriter.Flush()
		_ = tabWriter.Flush()

		return nil
	},
}
