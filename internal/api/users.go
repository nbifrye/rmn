package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

type usersResponse struct {
	Users      []User `json:"users"`
	TotalCount int    `json:"total_count"`
	Offset     int    `json:"offset"`
	Limit      int    `json:"limit"`
}

type userResponse struct {
	User User `json:"user"`
}

func (c *Client) ListUsers(ctx context.Context, params UserListParams) ([]User, int, error) {
	q := url.Values{}
	if params.Status > 0 {
		q.Set("status", strconv.Itoa(params.Status))
	}
	if params.Name != "" {
		q.Set("name", params.Name)
	}
	if params.GroupID > 0 {
		q.Set("group_id", strconv.Itoa(params.GroupID))
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset > 0 {
		q.Set("offset", strconv.Itoa(params.Offset))
	}

	var resp usersResponse
	if err := c.Get(ctx, "/users.json", q, &resp); err != nil {
		return nil, 0, err
	}
	return resp.Users, resp.TotalCount, nil
}

func (c *Client) GetUser(ctx context.Context, id int) (*User, error) {
	var resp userResponse
	if err := c.Get(ctx, fmt.Sprintf("/users/%d.json", id), nil, &resp); err != nil {
		return nil, err
	}
	return &resp.User, nil
}

func (c *Client) GetCurrentUser(ctx context.Context) (*User, error) {
	var resp userResponse
	if err := c.Get(ctx, "/users/current.json", nil, &resp); err != nil {
		return nil, err
	}
	return &resp.User, nil
}
