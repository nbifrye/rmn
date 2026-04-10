package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

type timeEntriesResponse struct {
	TimeEntries []TimeEntry `json:"time_entries"`
	TotalCount  int         `json:"total_count"`
	Offset      int         `json:"offset"`
	Limit       int         `json:"limit"`
}

type timeEntryResponse struct {
	TimeEntry TimeEntry `json:"time_entry"`
}

type timeEntryCreateRequest struct {
	TimeEntry TimeEntryCreateParams `json:"time_entry"`
}

type timeEntryUpdateRequest struct {
	TimeEntry TimeEntryUpdateParams `json:"time_entry"`
}

func (c *Client) ListTimeEntries(ctx context.Context, params TimeEntryListParams) ([]TimeEntry, int, error) {
	q := url.Values{}
	if params.ProjectID != "" {
		q.Set("project_id", params.ProjectID)
	}
	if params.IssueID > 0 {
		q.Set("issue_id", strconv.Itoa(params.IssueID))
	}
	if params.UserID > 0 {
		q.Set("user_id", strconv.Itoa(params.UserID))
	}
	if params.SpentOn != "" {
		q.Set("spent_on", params.SpentOn)
	}
	if params.From != "" {
		q.Set("from", params.From)
	}
	if params.To != "" {
		q.Set("to", params.To)
	}
	if params.ActivityID > 0 {
		q.Set("activity_id", strconv.Itoa(params.ActivityID))
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset > 0 {
		q.Set("offset", strconv.Itoa(params.Offset))
	}

	var resp timeEntriesResponse
	if err := c.Get(ctx, "/time_entries.json", q, &resp); err != nil {
		return nil, 0, err
	}
	return resp.TimeEntries, resp.TotalCount, nil
}

func (c *Client) GetTimeEntry(ctx context.Context, id int) (*TimeEntry, error) {
	var resp timeEntryResponse
	if err := c.Get(ctx, fmt.Sprintf("/time_entries/%d.json", id), nil, &resp); err != nil {
		return nil, err
	}
	return &resp.TimeEntry, nil
}

func (c *Client) CreateTimeEntry(ctx context.Context, params TimeEntryCreateParams) (*TimeEntry, error) {
	var resp timeEntryResponse
	if err := c.Post(ctx, "/time_entries.json", timeEntryCreateRequest{TimeEntry: params}, &resp); err != nil {
		return nil, err
	}
	return &resp.TimeEntry, nil
}

func (c *Client) UpdateTimeEntry(ctx context.Context, id int, params TimeEntryUpdateParams) error {
	return c.Put(ctx, fmt.Sprintf("/time_entries/%d.json", id), timeEntryUpdateRequest{TimeEntry: params})
}

func (c *Client) DeleteTimeEntry(ctx context.Context, id int) error {
	return c.Delete(ctx, fmt.Sprintf("/time_entries/%d.json", id))
}
