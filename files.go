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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
)

const filePath = "/api/v2/file"

// UploadFile a file to the file-manager.
// path must start with a slash (/)
func (cli *Client) UploadFile(ctx context.Context, path, name string, reader io.Reader) error {
	return cli.UploadFileReplace(ctx, path, name, false, reader)
}

// UploadFileReplace a file to the file-manager if replace is set to true
func (cli *Client) UploadFileReplace(ctx context.Context, path, name string, replace bool, reader io.Reader) error {
	buf := new(bytes.Buffer)
	multipartWriter := multipart.NewWriter(buf)

	partHeaders := make(textproto.MIMEHeader)
	partHeaders.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, name))
	partHeaders.Set("Content-Type", "")

	part, err := multipartWriter.CreatePart(partHeaders)
	if err != nil {
		return err
	}

	if _, err = io.Copy(part, reader); err != nil {
		return err
	}

	if part, err = multipartWriter.CreateFormField("path"); err != nil {
		return err
	}

	if _, err = part.Write([]byte(path)); err != nil {
		return err
	}

	if part, err = multipartWriter.CreateFormField("replace"); err != nil {
		return err
	}

	if _, err = part.Write([]byte(fmt.Sprintf("%t", replace))); err != nil {
		return err
	}

	if err = multipartWriter.Close(); err != nil {
		return err
	}

	requestURL := cli.baseURL + filePath

	var request *http.Request
	if request, err = http.NewRequestWithContext(ctx, http.MethodPost, requestURL, buf); err != nil {
		return err
	}

	request.Header.Add("Content-Type", multipartWriter.FormDataContentType())
	request.Header.Set("Authorization", "Bearer "+cli.authToken)

	var response *http.Response
	if response, err = cli.DoWithRefresh(request); err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

const fileMetadataPath = "/api/v2/file/metadata"

// GetFileMetadata returns the metadata for the given file.
func (cli *Client) GetFileMetadata(ctx context.Context, filePath string) (*File, error) {
	path := fileMetadataPath + "?path=" + filePath

	file := new(File)

	if err := cli.Call(ctx, http.MethodGet, path, nil, file); err != nil {
		return nil, err
	}

	return file, nil
}

type renameRequest struct {
	// Path to the file to rename.
	Path string `json:"path"`

	// New name of the file. If set, the NewPath field should be empty.
	Name string `json:"name"`

	// NewPath of the file. If set, the Name field should be empty.
	NewPath string `json:"newPath"`
}

// RenameFile in the file-manager.
func (cli *Client) RenameFile(ctx context.Context, path, name, newPath string) error {
	req := &renameRequest{
		Path:    path,
		Name:    name,
		NewPath: newPath,
	}

	return cli.Call(ctx, http.MethodPatch, filePath, req, nil)
}

// DeleteFile deletes a file.
func (cli *Client) DeleteFile(ctx context.Context, name string) error {
	path := filePath

	body := map[string]string{
		"path": name,
	}

	return cli.Call(ctx, http.MethodDelete, path, body, nil)
}

// DownloadFile from the file-manager.
// Returns a ReadCloser that must be closed after use.
func (cli *Client) DownloadFile(ctx context.Context, name string) (io.ReadCloser, error) {
	requestURL := cli.baseURL + filePath + "?path=" + name

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+cli.authToken)

	var response *http.Response

	if response, err = cli.DoWithRefresh(request); err != nil {
		return nil, err
	}

	return response.Body, nil
}

const fileCreateDirPath = "/api/v2/file/createdir"

type createDirectoryRequest struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

// CreateDirectory creates a directory with the given name in the given path.
func (cli *Client) CreateDirectory(ctx context.Context, path, name string) error {
	req := &createDirectoryRequest{
		Path: path,
		Name: name,
	}

	return cli.Call(ctx, http.MethodPost, fileCreateDirPath, req, nil)
}

// File represents a file in the file-manager.
type File struct {
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Mime      string `json:"mime"`
	Size      int    `json:"size"`
	Created   int64  `json:"created"`
	IsDir     bool   `json:"is_dir"`
	Path      string `json:"path"`
	Digest    string `json:"digest"`
	User      struct {
		ID string `json:"user_id"`
	} `json:"user"`
}

// FilesListResponse is a response from the file-manager containing a list of files.
type FilesListResponse struct {
	Items []File `json:"items"`
	Total int    `json:"total"`
}

const fileListPath = "/api/v2/files"

// ListSearch defines search parameters for ListQuery.
type ListSearch struct {
	// Name - file name to search for (partial-match).
	Name string `json:"name"`
	// Path - file path to search for (partial-match).
	Path string `json:"path"`
}

const (
	// SortDirectionAsc returns the files sorted in ascending order.
	SortDirectionAsc = "asc"

	// SortDirectionDesc returns the files sorted in descending order.
	SortDirectionDesc = "desc"
)

// ListQuery defines query parameters for FileList.
type ListQuery struct {
	// Path defines path to the directory to list files in.
	// If path is empty and:
	//	- search is empty - show files from the root directory
	//  - search is not empty - search in ALL files
	// Else, if path is not empty and:
	//	- search is empty - show files from that directory
	//  - search is not empty - search in that directory only
	Path string

	Search ListSearch

	// SortField defines field used to sort, 'name' by default.
	// Supported sort fields:
	// - name
	// - size
	// - extension
	// - created
	SortField string

	// SortDirection defines sort direction, 'asc' by default.
	SortDirection string

	// ItemsPerPage defines maximum number of records in result, default 50, max 1000
	ItemsPerPage int

	// Offset defines offset of the first record in result, default 0
	Offset int
}

// String returns string representation of CommitListQuery which can be used as query string.
func (q ListQuery) String() (string, error) {
	values := make(url.Values)

	if q.Path != "" {
		values.Set("path", q.Path)
	}

	searchQueryJSON, err := json.Marshal(q.Search)
	if err != nil {
		return "", err
	}

	values.Set("search", string(searchQueryJSON))

	if q.SortField != "" {
		values.Set("sort_field", q.SortField)
	}

	if q.SortDirection != "" {
		values.Set("sort_direction", q.SortDirection)
	}

	if q.ItemsPerPage != 0 {
		values.Set("items_per_page", fmt.Sprintf("%d", q.ItemsPerPage))
	}

	if q.Offset != 0 {
		values.Set("offset", fmt.Sprintf("%d", q.Offset))
	}

	return values.Encode(), nil
}

// ListFiles returns a list of files in the file-manager.
func (cli *Client) ListFiles(ctx context.Context, query ListQuery) (*FilesListResponse, error) {
	response := new(FilesListResponse)

	requestQueryParameters, err := query.String()
	if err != nil {
		return nil, fmt.Errorf("failed to encode query: %w", err)
	}

	requestPath := fileListPath + "?" + requestQueryParameters

	if err = cli.Call(ctx, http.MethodGet, requestPath, nil, response); err != nil {
		return nil, err
	}

	return response, nil
}

const fileStatsPath = "/api/v2/file/stats"

// Stats represents the statistics of file manager for a company.
type Stats struct {
	// CountFiles is the number of files in the file-manager.
	CountFiles int64 `json:"count_files"`

	// Quota is the allocated storage space in MB.
	Quota int64 `json:"quota"`

	// Used is the used storage space in MB.
	Used float64 `json:"used"`
}

// FileStats returns file-manager statistics.
func (cli *Client) FileStats(ctx context.Context) (*Stats, error) {
	stats := new(Stats)

	if err := cli.Call(ctx, http.MethodGet, fileStatsPath, nil, stats); err != nil {
		return nil, err
	}

	return stats, nil
}
