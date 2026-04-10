package api

import "context"

type statusesResponse struct {
	Statuses []IssueStatus `json:"issue_statuses"`
}

func (c *Client) ListStatuses(ctx context.Context) ([]IssueStatus, error) {
	var resp statusesResponse
	if err := c.Get(ctx, "/issue_statuses.json", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Statuses, nil
}
