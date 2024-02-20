package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"go.qbee.io/client/types"
)

// Commit represents a single commit in the system.
type Commit struct {
	// ID is the unique identifier of the commit.
	ID string `json:"-"`

	// SHA is the pseudo-digest of the commit.
	SHA string `json:"sha"`

	// Message is the commit message.
	Message string `json:"message"`

	// Type is the type of the commit.
	Type string `json:"type"`

	// Labels is the list of labels of the commit.
	Labels []string `json:"labels"`

	// Changes is the list of changes' SHA that are part of the commit.
	Changes []string `json:"changes"`

	// UserID is the ID of the user who created the commit.
	UserID string `json:"user_id"`

	// User contains the user who created the commit.
	// This field is populated by the API handlers.
	User *UserBaseInfo `json:"user,omitempty"`

	// Created is the timestamp when the commit was created.
	Created types.Timestamp `json:"created"`

	// Updated is the timestamp when the commit was last updated.
	Updated types.Timestamp `json:"updated"`
}

// CommitRequest is a request to commit uncommitted changes.
type CommitRequest struct {
	// Action must be always set to "commit".
	Action string `json:"action"`

	// Message describing changes in the commit.
	Message string `json:"message"`
}

const commitAction = "commit"

// CommitConfiguration commits uncommitted changes with provided message.
func (cli *Client) CommitConfiguration(ctx context.Context, message string) (*Commit, error) {
	const path = "/api/v2/commit"

	request := CommitRequest{
		Action:  commitAction,
		Message: message,
	}

	auditCommit := new(Commit)

	err := cli.Call(ctx, http.MethodPost, path, request, auditCommit)
	if err != nil {
		return nil, err
	}

	return auditCommit, nil
}

// CommitExtended is a commit with complete changes objects.
type CommitExtended struct {
	Commit

	// Changes is the list of changes' SHA that are part of the commit.
	Changes []Change `json:"changes"`
}

// GetCommit returns a single commit by its SHA.
func (cli *Client) GetCommit(ctx context.Context, sha string) (*CommitExtended, error) {
	path := "/api/v2/commit/" + sha

	commit := new(CommitExtended)

	if err := cli.Call(ctx, http.MethodGet, path, nil, commit); err != nil {
		return nil, err
	}

	return commit, nil
}

// CommitListSearch defines search parameters for CommitListQuery.
type CommitListSearch struct {
	// CommitSHA - full string or substring with at least 6 characters.
	CommitSHA string `json:"commit_sha,omitempty"`

	// Message - commit message to search for (substring match).
	Message string `json:"message,omitempty"`

	// UserID - user ID to search for (exact match).
	UserID string `json:"user_id,omitempty"`

	// Committer - name or surname of the user (substring match)
	Committer string `json:"committer,omitempty"`

	// ChangeSHA - full string or substring with at least 6 characters
	ChangeSHA string `json:"change_sha,omitempty"`

	// LabelsInclude - array of labels which MUST be present (exact match)
	LabelsInclude []string `json:"labels_include,omitempty"`

	// LabelsExclude - array of labels which MUST NOT be present (exact match)
	LabelsExclude []string `json:"labels_exclude,omitempty"`

	// StartDate - start date UTC, format: YYYY-MM-DD hh:ii:ss
	StartDate string `json:"start_date,omitempty"`

	// EndDate - end date UTC, format: YYYY-MM-DD hh:ii:ss
	EndDate string `json:"end_date,omitempty"`
}

// CommitListQuery defines query parameters for ChangeList.
type CommitListQuery struct {
	Search CommitListSearch

	// SortField defines field used to sort, 'created' by default.
	// Supported sort fields:
	// - user_id
	// - type
	// - message
	// - created
	Sort string

	// SortDirection defines sort direction, 'desc' by default.
	SortDirection string

	// Limit defines maximum number of records in result, default 30, max 1000
	Limit int

	// Offset defines offset of the first record in result, default 0
	Offset int
}

// String returns string representation of CommitListQuery which can be used as query string.
func (q CommitListQuery) String() (string, error) {
	values := make(url.Values)

	searchQueryJSON, err := json.Marshal(q.Search)
	if err != nil {
		return "", err
	}

	values.Set("search", string(searchQueryJSON))

	if q.Sort != "" {
		values.Set("sort_field", q.Sort)
	}

	if q.SortDirection != "" {
		values.Set("sort_direction", q.SortDirection)
	}

	if q.Limit != 0 {
		values.Set("items_per_page", fmt.Sprintf("%d", q.Limit))
	}

	if q.Offset != 0 {
		values.Set("offset", fmt.Sprintf("%d", q.Offset))
	}

	return values.Encode(), nil
}

// CommitsList represents a slice of commits matched by the query.
// As well as the total number of commits matched by the query.
type CommitsList struct {
	Commits []*CommitExtended `json:"items"`
	Total   int               `json:"total"`
}

// ListCommits returns a list of commits based on provided query.
func (cli *Client) ListCommits(ctx context.Context, query CommitListQuery) (*CommitsList, error) {
	queryString, err := query.String()
	if err != nil {
		return nil, err
	}

	path := "/api/v2/commitlist?" + queryString

	response := new(CommitsList)

	if err = cli.Call(ctx, http.MethodGet, path, nil, response); err != nil {
		return nil, err
	}

	return response, nil
}
