package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type projectsResponse struct {
	Projects   []Project `json:"projects"`
	TotalCount int       `json:"total_count"`
	Offset     int       `json:"offset"`
	Limit      int       `json:"limit"`
}

type projectResponse struct {
	Project Project `json:"project"`
}

type projectCreateRequest struct {
	Project ProjectCreateParams `json:"project"`
}

type projectUpdateRequest struct {
	Project ProjectUpdateParams `json:"project"`
}

func (c *Client) ListProjects(ctx context.Context, params ProjectListParams) ([]Project, int, error) {
	q := url.Values{}
	if params.Status != "" {
		q.Set("status", params.Status)
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset > 0 {
		q.Set("offset", strconv.Itoa(params.Offset))
	}

	var resp projectsResponse
	if err := c.Get(ctx, "/projects.json", q, &resp); err != nil {
		return nil, 0, err
	}
	return resp.Projects, resp.TotalCount, nil
}

func (c *Client) GetProject(ctx context.Context, id string, include []string) (*Project, error) {
	var q url.Values
	if len(include) > 0 {
		q = url.Values{}
		q.Set("include", strings.Join(include, ","))
	}
	var resp projectResponse
	if err := c.Get(ctx, fmt.Sprintf("/projects/%s.json", id), q, &resp); err != nil {
		return nil, err
	}
	return &resp.Project, nil
}

func (c *Client) CreateProject(ctx context.Context, params ProjectCreateParams) (*Project, error) {
	var resp projectResponse
	if err := c.Post(ctx, "/projects.json", projectCreateRequest{Project: params}, &resp); err != nil {
		return nil, err
	}
	return &resp.Project, nil
}

func (c *Client) UpdateProject(ctx context.Context, id string, params ProjectUpdateParams) error {
	return c.Put(ctx, fmt.Sprintf("/projects/%s.json", id), projectUpdateRequest{Project: params})
}

func (c *Client) ArchiveProject(ctx context.Context, id string) error {
	return c.Put(ctx, fmt.Sprintf("/projects/%s/archive.json", id), nil)
}

func (c *Client) UnarchiveProject(ctx context.Context, id string) error {
	return c.Put(ctx, fmt.Sprintf("/projects/%s/unarchive.json", id), nil)
}

func (c *Client) DeleteProject(ctx context.Context, id string) error {
	return c.Delete(ctx, fmt.Sprintf("/projects/%s.json", id))
}
