// Copyright 2023 qbee.io
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

package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/qbee-io/qbee-cli/client"
)

const (
	fileManagerSourceOption      = "source"
	fileManagerDestinationOption = "destination"
)

var filemanagerCommand = Command{
	Description: "Synchronize a local directory with the filemanager",
	SubCommands: map[string]Command{
		"sync":  fileManagerSyncCommand,
		"purge": fileManagerPurgeCommand,
		"print": fileManagerPrintCommand,
		"diff": {
			Description: "Diff a local directory with the filemanager",
			Target: func(opts Options) error {
				fmt.Printf("Diffing local directory %s with filemanager directory %s\n", opts[fileManagerSourceOption], opts[fileManagerDestinationOption])
				return nil
			},
		},
	},
}

var fileManagerSyncCommand = Command{
	Description: "Synchronize a local directory with the filemanager",
	Options: []Option{
		{
			Name:     fileManagerSourceOption,
			Short:    "s",
			Help:     "Local source path",
			Required: true,
		},
		{
			Name:     fileManagerDestinationOption,
			Short:    "d",
			Help:     "Destination path in the filemanager",
			Required: true,
		},
	},
	Target: func(opts Options) error {
		ctx := context.Background()

		cli, err := client.LoginGetAuthenticatedClient(ctx)

		if err != nil {
			return err
		}
		fileManager := client.NewFileManager().WithClient(cli).WithDryRun(false).WithDelete(true)

		startSync := time.Now()
		if err := fileManager.Sync(ctx, opts[fileManagerSourceOption], opts[fileManagerDestinationOption]); err != nil {
			return err
		}
		syncTime := (time.Now().UnixNano() - startSync.UnixNano()) / (int64(time.Millisecond) / int64(time.Nanosecond))
		s := fileManager.GetStatistics()

		fmt.Printf("Sync results:\nBytes written: %d\nFiles uploaded: %d\nTime spent: %d millisecond(s)\nFiles deleted: %d\n", s.Bytes, s.Files, syncTime, s.DeletedFiles)

		return nil
	},
}

var fileManagerPrintCommand = Command{
	Description: "Print a files in the filemanager",
	Options: []Option{
		{
			Name:     fileManagerDestinationOption,
			Short:    "d",
			Help:     "Destination path in the filemanager",
			Required: true,
		},
	},
	Target: func(opts Options) error {
		return nil
	},
}

var fileManagerPurgeCommand = Command{
	Description: "Purge a directory in the filemanager",
	Options: []Option{
		{
			Name:     fileManagerDestinationOption,
			Short:    "d",
			Help:     "Destination path in the filemanager",
			Required: true,
		},
	},
	Target: func(opts Options) error {
		/*
			ctx := context.Background()

			fileManager, err := GetAuthFileManager()
			if err != nil {
				return err
			}
		*/
		fmt.Printf("Purging directory %s\n", opts[fileManagerDestinationOption])

		return nil
	},
}
