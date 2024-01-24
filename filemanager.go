// Copyright 2023 qbee.io
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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// Manager manages the sync operation.
type Manager struct {
	client      *Client
	del         bool
	dryrun      bool
	verbose     bool
	remoteFiles fileMap
	localFiles  fileMap
}

type fileMap struct {
	mutex sync.Mutex
	files map[string]*fileInfo
}

func (fm *fileMap) add(key string, file *fileInfo) {
	fm.mutex.Lock()
	defer fm.mutex.Unlock()
	fm.files[key] = file
}

func (fm *fileMap) get(key string) (*fileInfo, bool) {
	fm.mutex.Lock()
	defer fm.mutex.Unlock()
	if file, ok := fm.files[key]; ok {
		return file, true
	}
	return nil, false
}

// fileInfo represents a file in the filemanager.
type fileInfo struct {
	File
}

// NewFileManager returns a new filemanager.
func NewFileManager() *Manager {
	return &Manager{}
}

// WithClient sets the client
func (fm *Manager) WithClient(client *Client) *Manager {
	fm.client = client
	return fm
}

// WithDelete sets the delete flag
func (m *Manager) WithDelete(del bool) *Manager {
	m.del = del
	return m
}

// WithDryRun sets the dryrun flag
func (m *Manager) WithDryRun(dryrun bool) *Manager {
	m.dryrun = dryrun
	return m
}

func (m *Manager) WithVerbose(verbose bool) *Manager {
	m.verbose = verbose
	return m
}

// Sync synchronizes the local directory with the filemanager.
func (m *Manager) Sync(ctx context.Context, source, dest string) error {

	if err := m.snapshotRemote(ctx, dest); err != nil {
		return err
	}

	if err := m.snapshotLocal(source); err != nil {
		return err
	}

	updates, err := m.filterUploads()

	if err != nil {
		return err
	}

	for _, uploadFile := range updates.files {
		if uploadFile.IsDir {
			continue
		}

		if err := m.upload(ctx, uploadFile, source, dest); err != nil {
			return err
		}
	}

	if m.del {
		deletes := m.filterDeletes()
		for _, remotePath := range deletes {
			if err := m.deleteRemote(remotePath); err != nil {
				return err
			}
		}
	}

	return nil
}

// Purge deletes all files in the given path.
func (m *Manager) Purge(ctx context.Context, remotePath string) error {

	if err := m.snapshotRemote(ctx, remotePath); err != nil {
		return err
	}

	deletes := m.filterDeletes()
	for _, deletePath := range deletes {
		if err := m.deleteRemote(deletePath); err != nil {
			return err
		}
	}

	return nil
}

// List prints the files in the filemanager.
func (m *Manager) List(ctx context.Context, remotePath string) error {

	if err := m.snapshotRemote(ctx, remotePath); err != nil {
		return err
	}

	for _, remoteFile := range m.remoteFiles.files {
		fmt.Printf("%s\n", remoteFile.Path)
	}
	return nil
}

// filterDeletes returns the files that need to be deleted.
func (m *Manager) filterDeletes() []string {

	deletes := []string{}

	for remoteRelativeName, remoteFile := range m.remoteFiles.files {
		if _, ok := m.localFiles.get(remoteRelativeName); !ok {
			deletes = append(deletes, remoteFile.Path)
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(deletes)))
	return deletes
}

// filterUploads returns the files that need to be uploaded.
func (m *Manager) filterUploads() (*fileMap, error) {

	uploads := fileMap{
		files: make(map[string]*fileInfo),
	}

	for localRelativeName, localFile := range m.localFiles.files {
		if localFile.IsDir {
			continue
		}

		remoteFile, ok := m.remoteFiles.get(localRelativeName)
		if !ok {
			uploads.add(localRelativeName, localFile)
			continue
		}

		if localFile.Size != remoteFile.Size {
			uploads.add(localRelativeName, localFile)
			continue
		}

		if localDigest, err := getFileDigest(localFile.Path); err != nil {
			return nil, err
		} else if localDigest != remoteFile.Digest {
			uploads.add(localRelativeName, localFile)
		}
	}
	return &uploads, nil
}

// snapshotLocal returns a channel which receives the infos of the files under the given basePath.
func (m *Manager) snapshotLocal(path string) error {

	m.localFiles.files = make(map[string]*fileInfo)

	basePath := filepath.ToSlash(filepath.Clean(path))
	stat, err := os.Stat(basePath)

	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return fmt.Errorf("the path %s is not a directory", basePath)
	}

	err = filepath.Walk(basePath, func(p string, stat os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fileName, err := filepath.Rel(basePath, p)
		if err != nil {
			return err
		}

		fileInfo := fileInfo{
			File: File{
				Name:  fileName,
				Path:  p,
				Size:  int(stat.Size()),
				IsDir: stat.IsDir(),
			},
		}
		m.localFiles.add(fileName, &fileInfo)
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// listFileManagerFiles returns a channel which receives the infos of the files under the given basePath.
func (m *Manager) snapshotRemote(ctx context.Context, remotePath string) error {

	m.remoteFiles.files = make(map[string]*fileInfo)

	basePath := path.Clean(remotePath)

	query := ListQuery{
		ItemsPerPage: 1000,
		Offset:       0,
		Search: ListSearch{
			Path: fmt.Sprintf("^%s/.*", basePath),
		},
	}

	for {
		files, err := m.client.ListFiles(ctx, query)
		if err != nil {
			return err
		}
		for _, file := range files.Items {

			if err != nil {
				return err
			}

			if strings.HasPrefix(file.Path, basePath) && file.Path != basePath {
				remoteRelativeName := strings.TrimPrefix(file.Path, basePath+"/")
				m.remoteFiles.add(remoteRelativeName, &fileInfo{File: file})
			}
		}

		query.Offset += query.ItemsPerPage
		if query.Offset > files.Total {
			break
		}
	}

	return nil
}

// upload uploads the file to the filemanager.
func (m *Manager) upload(ctx context.Context, file *fileInfo, sourcePath, destPath string) error {

	// We do not upload directories
	if file.IsDir {
		return nil
	}

	baseName := filepath.Base(file.Name)
	dirName := filepath.ToSlash(filepath.Dir(file.Name))
	destinationPath := filepath.ToSlash(filepath.Join(destPath, dirName))

	fmt.Printf("Uploading %s to %s\n", file.Path, destinationPath)

	if m.dryrun {
		return nil
	}

	reader, err := os.Open(file.Path)
	if err != nil {
		return err
	}

	defer reader.Close()

	if err := m.client.UploadFileReplace(ctx, destinationPath, baseName, true, reader); err != nil {
		return err
	}

	return nil
}

// deleteRemote deletes the remote file.
func (m *Manager) deleteRemote(remotePath string) error {

	println("Deleting", remotePath)
	if m.dryrun {
		return nil
	}

	return m.client.DeleteFile(context.Background(), remotePath)
}

// getFileDigest returns the sha256 digest of the given file.
func getFileDigest(src string) (string, error) {
	fp, err := os.Open(src)
	if err != nil {
		return "", fmt.Errorf("error opening file %s: %w", src, err)
	}
	defer fp.Close()

	if _, err = fp.Stat(); err != nil {
		return "", fmt.Errorf("error getting file metadata %s: %w", src, err)
	}

	digest := sha256.New()
	if _, err = io.Copy(digest, fp); err != nil {
		return "", fmt.Errorf("error calculating file checksum %s: %w", src, err)
	}

	return hex.EncodeToString(digest.Sum(nil)), nil
}
