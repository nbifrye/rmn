package api

import (
	"context"
	"fmt"
)

type versionsResponse struct {
	Versions   []Version `json:"versions"`
	TotalCount int       `json:"total_count"`
}

type versionResponse struct {
	Version Version `json:"version"`
}

type versionCreateRequest struct {
	Version VersionCreateParams `json:"version"`
}

type versionUpdateRequest struct {
	Version VersionUpdateParams `json:"version"`
}

func (c *Client) ListVersions(ctx context.Context, projectID string) ([]Version, int, error) {
	var resp versionsResponse
	if err := c.Get(ctx, fmt.Sprintf("/projects/%s/versions.json", projectID), nil, &resp); err != nil {
		return nil, 0, err
	}
	return resp.Versions, resp.TotalCount, nil
}

func (c *Client) GetVersion(ctx context.Context, id int) (*Version, error) {
	var resp versionResponse
	if err := c.Get(ctx, fmt.Sprintf("/versions/%d.json", id), nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Version, nil
}

func (c *Client) CreateVersion(ctx context.Context, projectID string, params VersionCreateParams) (*Version, error) {
	var resp versionResponse
	if err := c.Post(ctx, fmt.Sprintf("/projects/%s/versions.json", projectID), versionCreateRequest{Version: params}, &resp); err != nil {
		return nil, err
	}
	return &resp.Version, nil
}

func (c *Client) UpdateVersion(ctx context.Context, id int, params VersionUpdateParams) error {
	return c.Put(ctx, fmt.Sprintf("/versions/%d.json", id), versionUpdateRequest{Version: params})
}

func (c *Client) DeleteVersion(ctx context.Context, id int) error {
	return c.Delete(ctx, fmt.Sprintf("/versions/%d.json", id))
}
