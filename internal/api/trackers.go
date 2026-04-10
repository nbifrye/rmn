package api

import "context"

type trackersResponse struct {
	Trackers []IdName `json:"trackers"`
}

func (c *Client) ListTrackers(ctx context.Context) ([]IdName, error) {
	var resp trackersResponse
	if err := c.Get(ctx, "/trackers.json", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Trackers, nil
}
