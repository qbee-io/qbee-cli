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
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type multiErr struct {
	mu  sync.Mutex
	err []error
}

func (e *multiErr) Append(err error) {
	e.mu.Lock()
	e.err = append(e.err, err)
	e.mu.Unlock()
}

func (e *multiErr) Len() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.err)
}

func (e *multiErr) ErrOrNil() error {
	if e.Len() > 0 {
		return e
	}
	return nil
}

func (e *multiErr) Error() string {
	var errMsgs []string
	for _, err := range e.err {
		errMsgs = append(errMsgs, err.Error())
	}
	return strings.Join(errMsgs, "\n")
}

// Manager manages the sync operation.
type Manager struct {
	client     *Client
	nJobs      int
	del        bool
	dryrun     bool
	statistics SyncStatistics
}

// SyncStatistics captures the sync statistics.
type SyncStatistics struct {
	Bytes        int
	Files        int64
	DeletedFiles int64
	mutex        sync.RWMutex
}

type operation int

const (
	opUpdate operation = iota
	opDelete
)

// fileInfo represents a file in the filemanager.
type fileInfo struct {
	File
	err            error
	existsInSource bool
	singleFile     bool
}

type fileOp struct {
	*fileInfo
	op operation
}

const DefaultParallel = 10

// NewFileManager returns a new filemanager.
func NewFileManager() *Manager {
	return &Manager{
		nJobs: DefaultParallel,
	}
}

// WithClient sets the client
func (fm *Manager) WithClient(client *Client) *Manager {
	fm.client = client
	return fm
}

// WithParallel sets the number of parallel jobs
func (m *Manager) WithParallel(n int) *Manager {
	m.nJobs = n
	return m
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

// Sync syncs the files between filemanager and local disks.
func (m *Manager) Sync(ctx context.Context, source, dest string) error {
	innerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	chJob := make(chan func())
	var wg sync.WaitGroup
	for i := 0; i < m.nJobs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range chJob {
				job()
			}
		}()
	}
	defer func() {
		close(chJob)
		wg.Wait()
	}()
	return m.syncLocalToFileManager(innerCtx, chJob, source, dest)
}

// GetStatistics returns the structure that contains the sync statistics
func (m *Manager) GetStatistics() SyncStatistics {
	m.statistics.mutex.Lock()
	defer m.statistics.mutex.Unlock()
	return SyncStatistics{Bytes: m.statistics.Bytes, Files: m.statistics.Files, DeletedFiles: m.statistics.DeletedFiles}
}

// syncLocalToFileManager syncs the files between filemanager and local disks.
func (m *Manager) syncLocalToFileManager(ctx context.Context, chJob chan func(), sourcePath, destPath string) error {
	wg := &sync.WaitGroup{}
	errs := &multiErr{}

	dirsToDelete := make(map[string]fileInfo, 0)

	for source := range filterFilesForSync(
		listLocalFiles(ctx, sourcePath), m.listFileManagerFiles(ctx, destPath, destPath), m.del,
	) {
		if source.op == opDelete && source.IsDir {
			// We are not deleting directories
			dirsToDelete[source.Name] = *source.fileInfo
			continue
		}

		wg.Add(1)
		source := source

		chJob <- func() {
			defer wg.Done()
			if source.err != nil {
				errs.Append(source.err)
				return
			}
			switch source.op {
			case opUpdate:
				if err := m.upload(source.fileInfo, sourcePath, destPath); err != nil {
					errs.Append(err)
				}
			case opDelete:
				if err := m.deleteRemote(source.fileInfo, destPath); err != nil {
					errs.Append(err)
				}
			}
		}
	}
	wg.Wait()

	// delete any non-existing directories at the destination, order paths alphabetically
	// to ensure that we delete the deepest directories first
	keys := make([]string, 0)
	for k := range dirsToDelete {
		keys = append(keys, k)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	for _, dir := range keys {
		dirToDelete := dirsToDelete[dir]
		if err := m.deleteRemote(&dirToDelete, destPath); err != nil {
			errs.Append(err)
		}
	}

	return errs.ErrOrNil()
}

// listFileManagerFiles returns a channel which receives the infos of the files under the given basePath.
func (m *Manager) listFileManagerFiles(ctx context.Context, basePath, path string) chan *fileInfo {
	c := make(chan *fileInfo, 50000) // TODO: revisit this buffer size later

	go func() {
		defer close(c)

		if _, err := m.client.GetFileMetadata(ctx, path); err != nil {
			if apiError := make(Error); errors.As(err, &apiError) {
				if errorMessage, ok := apiError["error"].(map[string]any); ok {
					if errorMessage["code"].(float64) != http.StatusNotFound {
						sendErrorInfoToChannel(ctx, c, err)
						return
					}
				}
			}
		}
		m.listFileManagerFilesRecursively(ctx, c, basePath, basePath)
	}()

	return c
}

// listFileManagerFilesRecursively lists the files under the given path recursively.
func (m *Manager) listFileManagerFilesRecursively(ctx context.Context, c chan *fileInfo, basePath, path string) {

	absoluteBasePath := filepath.Clean(basePath)

	query := ListQuery{
		ItemsPerPage: 1000,
		Offset:       0,
		Search: ListSearch{
			Path: fmt.Sprintf("^%s/.*", absoluteBasePath),
		},
	}

	for {
		files, err := m.client.ListFiles(ctx, query)
		if err != nil {
			sendErrorInfoToChannel(ctx, c, err)
			return
		}
		for _, file := range files.Items {

			if strings.HasPrefix(file.Path, absoluteBasePath) && file.Path != absoluteBasePath {
				fi := fileInfo{File: file}
				sendFileInfoToChannel(ctx, c, basePath, file.Path, fi, false)
			}
		}

		query.Offset += query.ItemsPerPage
		if query.Offset > files.Total {
			break
		}
	}
}

// upload uploads the file to the filemanager.
func (m *Manager) upload(file *fileInfo, sourcePath, destPath string) error {

	if file.IsDir {
		return nil
	}

	sourceFilename := filepath.Join(sourcePath, file.Name)
	baseName := filepath.Base(file.Name)
	dirName := filepath.ToSlash(filepath.Dir(file.Name))
	destinationPath := filepath.ToSlash(filepath.Join(destPath, dirName))

	if file.singleFile {
		sourceFilename = sourcePath
		baseName = filepath.Base(file.Path)
		destinationPath = destPath
	}

	fmt.Printf("Uploading %s to %s\n", sourceFilename, destinationPath)
	if m.dryrun {
		return nil
	}

	reader, err := os.Open(sourceFilename)
	if err != nil {
		return err
	}

	defer reader.Close()

	if err := m.client.UploadFileReplace(context.Background(), destinationPath, baseName, true, reader); err != nil {
		return err
	}
	m.updateFileTransferStatistics(file.File.Size)

	return nil
}

// deleteRemote deletes the remote file.
func (m *Manager) deleteRemote(file *fileInfo, destPath string) error {

	deletePath := filepath.Join(destPath, file.Name)
	println("Deleting", deletePath)
	if m.dryrun {
		return nil
	}
	err := m.client.DeleteFile(context.Background(), deletePath)
	if err != nil {
		return err
	}
	m.incrementDeletedFiles()
	return nil
}

// updateSyncStatistics updates the statistics of the amount of bytes transferred for one file
func (m *Manager) updateFileTransferStatistics(written int) {
	m.statistics.mutex.Lock()
	defer m.statistics.mutex.Unlock()
	m.statistics.Files++
	m.statistics.Bytes += written
}

// incrementDeletedFiles increments the counter used to capture the number of remote files deleted during the synchronization process
func (m *Manager) incrementDeletedFiles() {
	m.statistics.mutex.Lock()
	defer m.statistics.mutex.Unlock()
	m.statistics.DeletedFiles++
}

// listLocalFiles returns a channel which receives the infos of the files under the given basePath.
// basePath have to be absolute path.
func listLocalFiles(ctx context.Context, basePath string) chan *fileInfo {
	c := make(chan *fileInfo)

	basePath = filepath.ToSlash(basePath)

	go func() {
		defer close(c)

		stat, err := os.Stat(basePath)
		if os.IsNotExist(err) {
			// The path doesn't exist.
			// Returns and closes the channel without sending any.
			return
		} else if err != nil {
			sendErrorInfoToChannel(ctx, c, err)
			return
		}

		fileObject := fileInfo{
			File: File{
				Path:  basePath,
				Size:  int(stat.Size()),
				IsDir: stat.IsDir(),
			},
			singleFile: true,
		}

		if !stat.IsDir() {

			sendFileInfoToChannel(ctx, c, basePath, basePath, fileObject, true)
			return
		}

		sendFileInfoToChannel(ctx, c, basePath, basePath, fileObject, false)

		err = filepath.Walk(basePath, func(path string, stat os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fileInfo := fileInfo{
				File: File{
					Name:  filepath.ToSlash(path),
					Path:  basePath,
					Size:  int(stat.Size()),
					IsDir: stat.IsDir(),
				},
			}
			sendFileInfoToChannel(ctx, c, basePath, path, fileInfo, false)
			return ctx.Err()
		})

		if err != nil {
			sendErrorInfoToChannel(ctx, c, err)
		}

	}()
	return c
}

// sendFileInfoToChannel sends the file info to the channel.
func sendFileInfoToChannel(ctx context.Context, c chan *fileInfo, basePath, path string, fi fileInfo, singleFile bool) {
	relPath, _ := filepath.Rel(basePath, path)
	fi.Name = filepath.ToSlash(relPath)
	select {
	case c <- &fi:
	case <-ctx.Done():
	}
}

// sendErrorInfoToChannel sends the error to the channel.
func sendErrorInfoToChannel(ctx context.Context, c chan *fileInfo, err error) {
	fi := &fileInfo{
		err: err,
	}
	select {
	case c <- fi:
	case <-ctx.Done():
	}
}

// filterFilesForSync filters the source files from the given destination files, and returns
// another channel which includes the files necessary to be synced.
func filterFilesForSync(sourceFileChan, destFileChan chan *fileInfo, del bool) chan *fileOp {
	c := make(chan *fileOp)

	destFiles, err := fileInfoChanToMap(destFileChan)

	go func() {
		defer close(c)
		if err != nil {
			c <- &fileOp{fileInfo: &fileInfo{err: err}}
			return
		}
		for sourceInfo := range sourceFileChan {
			destInfo, ok := destFiles[sourceInfo.Name]
			// source is necessary to sync if
			// 1. The dest doesn't exist
			// 2. The dest doesn't have the same size as the source
			if !ok || sourceInfo.Size != destInfo.Size {
				c <- &fileOp{fileInfo: sourceInfo}
			} else {

				sourcePath := filepath.Join(sourceInfo.Path, sourceInfo.Name)
				sourceDigest, err := getFileDigest(sourcePath)
				if err != nil {
					c <- &fileOp{fileInfo: &fileInfo{err: err}}
					return
				}
				if sourceDigest != destInfo.Digest {
					c <- &fileOp{fileInfo: sourceInfo}
				}
			}

			if ok {
				destInfo.existsInSource = true
			}
		}
		if del {
			for _, destInfo := range destFiles {
				if !destInfo.existsInSource {
					// The source doesn't exist
					c <- &fileOp{fileInfo: destInfo, op: opDelete}
				}
			}
		}
	}()

	return c
}

// fileInfoChanToMap accumulates the fileInfos from the given channel and returns a map.
// It returns an error if the channel contains an error.
func fileInfoChanToMap(files chan *fileInfo) (map[string]*fileInfo, error) {
	result := make(map[string]*fileInfo)

	for file := range files {
		if file.err != nil {
			return nil, file.err
		}
		result[file.Name] = file
	}
	return result, nil
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
