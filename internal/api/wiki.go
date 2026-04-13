package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

type wikiPagesResponse struct {
	WikiPages []WikiPage `json:"wiki_pages"`
}

type wikiPageResponse struct {
	WikiPage WikiPageDetail `json:"wiki_page"`
}

type wikiPageWriteRequest struct {
	WikiPage WikiPageCreateParams `json:"wiki_page"`
}

type wikiPageUpdateRequest struct {
	WikiPage WikiPageUpdateParams `json:"wiki_page"`
}

func (c *Client) ListWikiPages(ctx context.Context, projectID string) ([]WikiPage, error) {
	var resp wikiPagesResponse
	if err := c.Get(ctx, fmt.Sprintf("/projects/%s/wiki/index.json", projectID), nil, &resp); err != nil {
		return nil, err
	}
	return resp.WikiPages, nil
}

func (c *Client) GetWikiPage(ctx context.Context, projectID string, title string, version int) (*WikiPageDetail, error) {
	path := fmt.Sprintf("/projects/%s/wiki/%s.json", projectID, url.PathEscape(title))
	var q url.Values
	if version > 0 {
		q = url.Values{}
		q.Set("version", strconv.Itoa(version))
	}
	var resp wikiPageResponse
	if err := c.Get(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return &resp.WikiPage, nil
}

func (c *Client) CreateWikiPage(ctx context.Context, projectID string, title string, params WikiPageCreateParams) (*WikiPageDetail, error) {
	var resp wikiPageResponse
	if err := c.PutWithResult(ctx, fmt.Sprintf("/projects/%s/wiki/%s.json", projectID, url.PathEscape(title)), wikiPageWriteRequest{WikiPage: params}, &resp); err != nil {
		return nil, err
	}
	return &resp.WikiPage, nil
}

func (c *Client) UpdateWikiPage(ctx context.Context, projectID string, title string, params WikiPageUpdateParams) error {
	return c.Put(ctx, fmt.Sprintf("/projects/%s/wiki/%s.json", projectID, url.PathEscape(title)), wikiPageUpdateRequest{WikiPage: params})
}

func (c *Client) DeleteWikiPage(ctx context.Context, projectID string, title string) error {
	return c.Delete(ctx, fmt.Sprintf("/projects/%s/wiki/%s.json", projectID, url.PathEscape(title)))
}
