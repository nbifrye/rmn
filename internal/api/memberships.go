package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

type membershipsResponse struct {
	Memberships []Membership `json:"memberships"`
	TotalCount  int          `json:"total_count"`
	Offset      int          `json:"offset"`
	Limit       int          `json:"limit"`
}

type membershipResponse struct {
	Membership Membership `json:"membership"`
}

type membershipCreateRequest struct {
	Membership MembershipCreateParams `json:"membership"`
}

type membershipUpdateRequest struct {
	Membership MembershipUpdateParams `json:"membership"`
}

func (c *Client) ListMemberships(ctx context.Context, projectID string, params MembershipListParams) ([]Membership, int, error) {
	q := url.Values{}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset > 0 {
		q.Set("offset", strconv.Itoa(params.Offset))
	}

	var resp membershipsResponse
	if err := c.Get(ctx, fmt.Sprintf("/projects/%s/memberships.json", projectID), q, &resp); err != nil {
		return nil, 0, err
	}
	return resp.Memberships, resp.TotalCount, nil
}

func (c *Client) GetMembership(ctx context.Context, id int) (*Membership, error) {
	var resp membershipResponse
	if err := c.Get(ctx, fmt.Sprintf("/memberships/%d.json", id), nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Membership, nil
}

func (c *Client) CreateMembership(ctx context.Context, projectID string, params MembershipCreateParams) (*Membership, error) {
	var resp membershipResponse
	if err := c.Post(ctx, fmt.Sprintf("/projects/%s/memberships.json", projectID), membershipCreateRequest{Membership: params}, &resp); err != nil {
		return nil, err
	}
	return &resp.Membership, nil
}

func (c *Client) UpdateMembership(ctx context.Context, id int, params MembershipUpdateParams) error {
	return c.Put(ctx, fmt.Sprintf("/memberships/%d.json", id), membershipUpdateRequest{Membership: params})
}

func (c *Client) DeleteMembership(ctx context.Context, id int) error {
	return c.Delete(ctx, fmt.Sprintf("/memberships/%d.json", id))
}
