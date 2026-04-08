package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

type issuesResponse struct {
	Issues     []Issue `json:"issues"`
	TotalCount int     `json:"total_count"`
	Offset     int     `json:"offset"`
	Limit      int     `json:"limit"`
}

type issueResponse struct {
	Issue Issue `json:"issue"`
}

type issueCreateRequest struct {
	Issue IssueCreateParams `json:"issue"`
}

type issueUpdateRequest struct {
	Issue IssueUpdateParams `json:"issue"`
}

func (c *Client) ListIssues(ctx context.Context, params IssueListParams) ([]Issue, int, error) {
	q := url.Values{}
	if params.ProjectID != "" {
		q.Set("project_id", params.ProjectID)
	}
	if params.StatusID != "" {
		q.Set("status_id", params.StatusID)
	}
	if params.AssignedToID != "" {
		q.Set("assigned_to_id", params.AssignedToID)
	}
	if params.TrackerID != 0 {
		q.Set("tracker_id", strconv.Itoa(params.TrackerID))
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset > 0 {
		q.Set("offset", strconv.Itoa(params.Offset))
	}

	var resp issuesResponse
	if err := c.Get(ctx, "/issues.json", q, &resp); err != nil {
		return nil, 0, err
	}
	return resp.Issues, resp.TotalCount, nil
}

func (c *Client) GetIssue(ctx context.Context, id int) (*Issue, error) {
	var resp issueResponse
	if err := c.Get(ctx, fmt.Sprintf("/issues/%d.json", id), nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Issue, nil
}

func (c *Client) CreateIssue(ctx context.Context, params IssueCreateParams) (*Issue, error) {
	var resp issueResponse
	if err := c.Post(ctx, "/issues.json", issueCreateRequest{Issue: params}, &resp); err != nil {
		return nil, err
	}
	return &resp.Issue, nil
}

func (c *Client) UpdateIssue(ctx context.Context, id int, params IssueUpdateParams) error {
	return c.Put(ctx, fmt.Sprintf("/issues/%d.json", id), issueUpdateRequest{Issue: params})
}

func (c *Client) DeleteIssue(ctx context.Context, id int) error {
	return c.Delete(ctx, fmt.Sprintf("/issues/%d.json", id))
}
