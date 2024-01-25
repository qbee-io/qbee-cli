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
	"path"
	"path/filepath"
	"time"

	"go.qbee.io/client"
)

const (
	fileManagerSourceOption      = "source"
	fileManagerDestinationOption = "destination"
	fileManagerDryRynOption      = "dry-run"
	fileManagerDeleteOption      = "delete"
	fileManagerRecursiveOption   = "recursive"
)

var filemanagerCommand = Command{
	Description: "Filemanager commands	",
	SubCommands: map[string]Command{
		"sync": fileManagerSyncCommand,
		"rm":   fileManagerPurgeCommand,
		"list": fileManagerListCommand,
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

		fileManager := client.NewFileManager().
			WithClient(cli).
			WithDelete(opts[fileManagerDeleteOption] == "true").
			WithDryRun(opts[fileManagerDryRynOption] == "true")

		remotePath := path.Clean(opts[fileManagerDestinationOption])
		localPath := filepath.Clean(opts[fileManagerSourceOption])

		fmt.Printf("Syncing directory %s to %s\n", localPath, remotePath)

		startSync := time.Now()
		if err := fileManager.Sync(ctx, localPath, remotePath); err != nil {
			return err
		}

		fmt.Printf("Time spent: %s\n", time.Since(startSync).String())
		return nil
	},
}

var fileManagerListCommand = Command{
	Description: "List files in the filemanager",
	Options: []Option{
		{
			Name:     fileManagerDestinationOption,
			Short:    "d",
			Help:     "Destination path in the filemanager",
			Default:  "/",
			Required: false,
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

		remotePath := path.Clean(opts[fileManagerDestinationOption])

		fmt.Printf("Listing directory %s\n", remotePath)

		if err := fileManager.List(ctx, remotePath); err != nil {
			return err
		}
		return nil

	},
}

var fileManagerPurgeCommand = Command{
	Description: "Remove a path in the filemanager",
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
			Name:     fileManagerRecursiveOption,
			Help:     "Recursive. Delete all files in the directory",
			Required: false,
			Flag:     "true",
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
			WithClient(cli).
			WithDelete(true).
			WithDryRun(opts[fileManagerDryRynOption] == "true")

		remotePath := path.Clean(opts[fileManagerDestinationOption])

		// Do not delete the root directory
		if remotePath == "/" {
			return fmt.Errorf("cannot delete root directory")
		}

		fmt.Printf("Removing path %s\n", remotePath)

		startSync := time.Now()
		if err := fileManager.Remove(ctx, remotePath, opts[fileManagerRecursiveOption] == "true"); err != nil {
			return err
		}
		fmt.Printf("Time spent: %s\n", time.Since(startSync).String())
		return nil
	},
}
