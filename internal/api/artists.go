package api

import (
	"encoding/json"
	"fmt"
)

func (c *Client) GetArtistBriefInfo(id int) (*ArtistBriefInfo, error) {
	raw, err := c.get(fmt.Sprintf("/artists/%d/brief-info", id), nil)
	if err != nil {
		return nil, err
	}
	var info ArtistBriefInfo
	if err := json.Unmarshal(raw, &info); err != nil {
		return nil, err
	}
	return &info, nil
}
