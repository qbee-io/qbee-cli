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
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// FileManager provides implementation of common batch operations inside qbee's file manager.
type FileManager struct {
	// client is the FileManager client.
	client *Client
	// deleteMissing deletes files on the remote that are not present locally.
	deleteMissing bool
	// dryRun does not perform any changes.
	dryRun bool
	// excludes removes local paths from being considered
	excludes []regexp.Regexp
	// includes includes only local paths that match the pattern
	includes []regexp.Regexp
	// remoteFiles is a map of the remote files.
	remoteFiles map[string]File
	// localFiles is a map of the local files.
	localFiles map[string]File
}

// NewFileManager returns a new FileManager.
func NewFileManager() *FileManager {
	return &FileManager{
		remoteFiles: make(map[string]File),
		localFiles:  make(map[string]File),
	}
}

// WithClient sets the client
func (m *FileManager) WithClient(client *Client) *FileManager {
	m.client = client
	return m
}

// WithDelete sets the delete flag
func (m *FileManager) WithDelete(del bool) *FileManager {
	m.deleteMissing = del
	return m
}

// WithDryRun sets the dryrun flag
func (m *FileManager) WithDryRun(dryrun bool) *FileManager {
	m.dryRun = dryrun
	return m
}

// GetLocalSnapshot returns the snapshot of the local files.
func (m *FileManager) GetLocalSnapshot() map[string]File {
	return m.localFiles
}

// GetRemoteSnapshot returns the snapshot of the remote files.
func (m *FileManager) GetRemoteSnapshot() map[string]File {
	return m.remoteFiles
}

// WithExcludes sets the excludes
func (m *FileManager) WithExcludes(excludes string) *FileManager {
	m.excludes = preparePathMatchList(excludes)
	return m
}

// WithIncludes sets the includes
func (m *FileManager) WithIncludes(includes string) *FileManager {
	m.includes = preparePathMatchList(includes)
	return m
}

// preparePathMatchList prepares a list of regular expressions from a comma-separated list of patterns.
func preparePathMatchList(patterns string) []regexp.Regexp {
	patternsList := strings.Split(patterns, ",")

	pathMatchList := make([]regexp.Regexp, 0)

	for _, pathPattern := range patternsList {
		pathPattern = strings.TrimLeft(pathPattern, "/")
		pathPattern = regexp.QuoteMeta(pathPattern)

		pathRE := regexp.MustCompile(pathPattern)
		pathMatchList = append(pathMatchList, *pathRE)
	}

	return pathMatchList
}

// Sync synchronizes the local directory with the FileManager.
func (m *FileManager) Sync(ctx context.Context, source, dest string) error {

	if err := m.SnapshotLocal(source); err != nil {
		return err
	}

	if err := m.SnapshotRemote(ctx, dest); err != nil {
		return err
	}

	updates, err := m.filterUploads()

	if err != nil {
		return err
	}

	for _, uploadFile := range updates {
		if uploadFile.IsDir {
			continue
		}

		if err := m.upload(ctx, uploadFile, dest); err != nil {
			return err
		}
	}

	if m.deleteMissing {
		// remove all files that from remote that are present locally
		for remoteRelativeName := range m.remoteFiles {
			if _, ok := m.localFiles[remoteRelativeName]; ok {
				delete(m.remoteFiles, remoteRelativeName)
			}
		}
		// delete all remaining files
		return m.deleteRemoteRecursive()
	}

	return nil
}

// Remove all files under provided remotePath.
//
// If recursive flag is set to true and remotePath points to a directory,
// this method will remove all nested files and directories recursively.
// Otherwise, it will return an error when deleting a non-empty directory.
func (m *FileManager) Remove(ctx context.Context, remotePath string, recursive bool) error {

	// Add the top directory to the list of files to delete
	m.remoteFiles[remotePath] = File{
		Name:  remotePath,
		Path:  remotePath,
		IsDir: true,
	}

	if recursive {
		if err := m.SnapshotRemote(ctx, remotePath); err != nil {
			return err
		}
	}

	return m.deleteRemoteRecursive()
}

// List files under provided remotePath.
func (m *FileManager) List(ctx context.Context, remotePath string) error {

	if err := m.SnapshotRemote(ctx, remotePath); err != nil {
		return err
	}

	for _, remoteFile := range sortFileMap(m.remoteFiles, false) {
		fmt.Println(remoteFile.Path)
	}
	return nil
}

// deleteRemoteRecursive deletes all discovered remote files in the FileManager.
func (m *FileManager) deleteRemoteRecursive() error {
	deletes := sortFileMap(m.remoteFiles, true)

	for _, remoteFile := range deletes {
		if err := m.deleteRemote(remoteFile); err != nil {
			return err
		}
	}
	return nil
}

// sortFileMap returns a slice of files sorted by path
func sortFileMap(fileMap map[string]File, reverse bool) []File {
	files := []File{}
	for _, file := range fileMap {
		files = append(files, file)
	}
	sort.Slice(files, func(i, j int) bool {
		if reverse {
			return files[i].Path > files[j].Path
		}
		return files[i].Path < files[j].Path
	})
	return files
}

// filterUploads returns the files that need to be uploaded.
func (m *FileManager) filterUploads() (map[string]File, error) {

	uploads := make(map[string]File)

	for localRelativeName, localFile := range m.localFiles {
		if localFile.IsDir {
			continue
		}

		remoteFile, ok := m.remoteFiles[localRelativeName]
		if !ok {
			uploads[localRelativeName] = localFile
			continue
		}

		if localFile.Size != remoteFile.Size {
			uploads[localRelativeName] = localFile
			continue
		}

		if localDigest, err := getFileDigest(localFile.Path); err != nil {
			return nil, err
		} else if localDigest != remoteFile.Digest {
			uploads[localRelativeName] = localFile
		}
	}
	return uploads, nil
}

// SnapshotLocal populates the localFiles map with the files in the local directory.
func (m *FileManager) SnapshotLocal(localPath string) error {
	m.localFiles = make(map[string]File)

	stat, err := os.Stat(localPath)

	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return fmt.Errorf("path %s is not a directory", localPath)
	}

	err = filepath.Walk(localPath, func(path string, stat os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fileName, err := filepath.Rel(localPath, path)
		if err != nil {
			return err
		}

		if m.matchExcludes(fileName) && !m.matchIncludes(localPath, fileName) {
			return nil
		}

		fileName = filepath.ToSlash(fileName)
		m.localFiles[fileName] = File{
			Name:  fileName,
			Path:  path,
			Size:  int(stat.Size()),
			IsDir: stat.IsDir(),
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// matchExcludes determines if a path matches any of the exclude filters.
func (m *FileManager) matchExcludes(path string) bool {
	if len(m.excludes) == 0 {
		return false
	}

	return matchRElist(path, m.excludes)
}

// matchIncludes determines if a path matches any of the include filters.
func (m *FileManager) matchIncludes(basePath, path string) bool {
	if len(m.includes) == 0 {
		return false
	}

	if !matchRElist(path, m.includes) {
		return false
	}

	// add directories to the list of files
	dirNameRelative := filepath.Dir(path)
	subDirPaths := strings.Split(dirNameRelative, string(filepath.Separator))

	for i := range subDirPaths {
		subDirPath := filepath.Join(subDirPaths[:i+1]...)
		subDirPathKey := filepath.ToSlash(subDirPath)
		if _, ok := m.localFiles[subDirPathKey]; ok {
			continue
		}

		absolutePath := filepath.Join(basePath, subDirPath)
		fileInfo, err := os.Stat(absolutePath)

		// this should never happen
		if os.IsNotExist(err) {
			return false
		}

		m.localFiles[subDirPathKey] = File{
			Name:  subDirPathKey,
			IsDir: true,
			Path:  absolutePath,
			Size:  int(fileInfo.Size()),
		}
	}
	return true
}

// matchRElist determines if a pattern matches a list of regular expressions.
func matchRElist(pattern string, list []regexp.Regexp) bool {
	if len(list) == 0 {
		return false
	}

	for _, re := range list {
		if re.MatchString(pattern) {
			return true
		}
	}
	return false
}

// SnapshotRemote populates the remoteFiles map with the files in the FileManager.
func (m *FileManager) SnapshotRemote(ctx context.Context, remotePath string) error {

	m.remoteFiles = make(map[string]File)

	searchPath := fmt.Sprintf("^%s/.*", remotePath)
	trimPrefix := remotePath + "/"

	// Search for the root directory
	if remotePath == "/" {
		searchPath = fmt.Sprintf("^%s.*", remotePath)
		trimPrefix = remotePath
	}

	query := ListQuery{
		ItemsPerPage: 1000,
		Offset:       0,
		Search: ListSearch{
			Path: searchPath,
		},
	}

	for {
		files, err := m.client.ListFiles(ctx, query)
		if err != nil {
			return err
		}
		for _, file := range files.Items {
			remoteRelativeName := strings.TrimPrefix(file.Path, trimPrefix)
			m.remoteFiles[remoteRelativeName] = file
		}

		query.Offset += query.ItemsPerPage
		if query.Offset > files.Total {
			break
		}
	}

	return nil
}

// upload uploads the file to the FileManager.
func (m *FileManager) upload(ctx context.Context, file File, destPath string) error {

	// We do not upload directories
	if file.IsDir {
		return nil
	}

	baseName := filepath.Base(file.Name)
	dirName := filepath.ToSlash(filepath.Dir(file.Name))
	destinationPath := filepath.ToSlash(filepath.Join(destPath, dirName))

	fmt.Printf("Uploading %s to %s\n", file.Path, destinationPath)

	if m.dryRun {
		return nil
	}

	reader, err := os.Open(file.Path)
	if err != nil {
		return err
	}
	defer func() { _ = reader.Close() }()

	if err := m.client.UploadFileReplace(ctx, destinationPath, baseName, true, reader); err != nil {
		return err
	}

	return nil
}

// deleteRemote deletes the remote file.
func (m *FileManager) deleteRemote(remoteFile File) error {

	fmt.Printf("Deleting %s\n", remoteFile.Path)
	if m.dryRun {
		return nil
	}

	return m.client.DeleteFile(context.Background(), remoteFile.Path)
}

// getFileDigest returns the sha256 digest of the given file.
func getFileDigest(src string) (string, error) {
	fp, err := os.Open(src)
	if err != nil {
		return "", fmt.Errorf("error opening file %s: %w", src, err)
	}
	defer func() { _ = fp.Close() }()

	if _, err = fp.Stat(); err != nil {
		return "", fmt.Errorf("error getting file metadata %s: %w", src, err)
	}

	digest := sha256.New()
	if _, err = io.Copy(digest, fp); err != nil {
		return "", fmt.Errorf("error calculating file checksum %s: %w", src, err)
	}

	return hex.EncodeToString(digest.Sum(nil)), nil
}

// UploadFile uploads a file to the FileManager.
func (m *FileManager) UploadFile(ctx context.Context, remotePath, localPath string, overwrite bool) error {
	reader, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer func() { _ = reader.Close() }()

	return m.client.UploadFileReplace(ctx, remotePath, localPath, overwrite, reader)
}

// DownloadFile downloads a file from the FileManager.
func (m *FileManager) DownloadFile(ctx context.Context, remotePath, localPath string) error {
	writer, err := os.Create(localPath)
	if err != nil {
		return err
	}

	defer func() { _ = writer.Close() }()

	reader, err := m.client.DownloadFile(ctx, remotePath)
	if err != nil {
		return err
	}

	if _, err := io.Copy(writer, reader); err != nil {
		return err
	}
	return nil
}
