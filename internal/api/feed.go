package api

import "encoding/json"

func (c *Client) GetFeed() (*FeedResponse, error) {
	raw, err := c.get("/feed", nil)
	if err != nil {
		return nil, err
	}
	var feed FeedResponse
	if err := json.Unmarshal(raw, &feed); err != nil {
		return nil, err
	}
	return &feed, nil
}
