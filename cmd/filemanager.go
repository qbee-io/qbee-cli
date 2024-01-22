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

package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.qbee.io/client"
)

const (
	fileManagerSourceOption      = "source"
	fileManagerDestinationOption = "destination"
	fileManagerDryRynOption      = "dry-run"
	fileManagerParalellOption    = "parallel"
	fileManagerDeleteOption      = "delete"
)

var filemanagerCommand = Command{
	Description: "Synchronize a local directory with the filemanager",
	SubCommands: map[string]Command{
		"sync":  fileManagerSyncCommand,
		"purge": fileManagerPurgeCommand,
		"print": fileManagerPrintCommand,
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
		{
			Name: fileManagerDryRynOption,
			Help: "Dry run. Do nothing, just print",
			Flag: "true",
		},
		{
			Name:    fileManagerParalellOption,
			Help:    "Number of parallel workers",
			Default: fmt.Sprintf("%d", client.DefaultParallel),
		},
		{
			Name: fileManagerDeleteOption,
			Help: "Delete files in the destination that are not in the source",
			Flag: "true",
		},
	},
	Target: func(opts Options) error {
		ctx := context.Background()

		cli, err := client.LoginGetAuthenticatedClient(ctx)

		if err != nil {
			return err
		}

		paralell, err := strconv.Atoi(opts[fileManagerParalellOption])
		if err != nil {
			return err
		}

		fileManager := client.
			NewFileManager().
			WithClient(cli).
			WithDelete(opts[fileManagerDeleteOption] == "true").
			WithDryRun(opts[fileManagerDryRynOption] == "true").
			WithParallel(paralell)

		fmt.Printf("Syncing directory %s to %s\n", opts[fileManagerSourceOption], opts[fileManagerDestinationOption])

		startSync := time.Now()
		if err := fileManager.Sync(ctx, opts[fileManagerSourceOption], opts[fileManagerDestinationOption]); err != nil {
			return err
		}
		syncTime := (time.Now().UnixNano() - startSync.UnixNano()) / (int64(time.Millisecond) / int64(time.Nanosecond))
		s := fileManager.GetStatistics()

		fmt.Printf("Sync results:\nBytes written: %d\n", s.Bytes)
		fmt.Printf("Files uploaded: %d\n", s.Files)
		fmt.Printf("Files deleted: %d\n", s.DeletedFiles)
		fmt.Printf("Directories deleted: %d\n", s.DeletedDirs)
		fmt.Printf("Time spent: %d millisecond(s)\n", syncTime)
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

		ctx := context.Background()

		cli, err := client.LoginGetAuthenticatedClient(ctx)

		if err != nil {
			return err
		}

		fileManager := client.
			NewFileManager().
			WithClient(cli)

		fmt.Printf("Printing directory %s\n", opts[fileManagerDestinationOption])

		if err := fileManager.Print(ctx, opts[fileManagerDestinationOption]); err != nil {
			return err
		}

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
		{
			Name:     fileManagerDryRynOption,
			Help:     "Dry run. Do nothing, just print",
			Required: false,
			Default:  "false",
			Flag:     "true",
		},
		{
			Name:     fileManagerParalellOption,
			Help:     "Number of parallel workers",
			Required: false,
			Default:  fmt.Sprintf("%d", client.DefaultParallel),
		},
	},
	Target: func(opts Options) error {

		ctx := context.Background()

		cli, err := client.LoginGetAuthenticatedClient(ctx)

		if err != nil {
			return err
		}

		paralell, err := strconv.Atoi(opts[fileManagerParalellOption])
		if err != nil {
			return err
		}

		fileManager := client.
			NewFileManager().
			WithClient(cli).
			WithDelete(true).
			WithDryRun(opts[fileManagerDryRynOption] == "true").
			WithParallel(paralell)

		fmt.Printf("Purging directory %s\n", opts[fileManagerDestinationOption])

		startSync := time.Now()
		if err := fileManager.Purge(ctx, opts[fileManagerDestinationOption]); err != nil {
			return err
		}
		syncTime := (time.Now().UnixNano() - startSync.UnixNano()) / (int64(time.Millisecond) / int64(time.Nanosecond))
		s := fileManager.GetStatistics()

		fmt.Printf("Purge results:\n")
		fmt.Printf("Files deleted: %d\n", s.DeletedFiles)
		fmt.Printf("Directories deleted: %d\n", s.DeletedDirs)
		fmt.Printf("Time spent: %d millisecond(s)\n", syncTime)
		return nil
	},
}
