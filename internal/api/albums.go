package api

import (
	"encoding/json"
	"fmt"
)

func (c *Client) GetAlbumWithTracks(id int) (*Album, error) {
	raw, err := c.get(fmt.Sprintf("/albums/%d/with-tracks", id), nil)
	if err != nil {
		return nil, err
	}
	var album Album
	if err := json.Unmarshal(raw, &album); err != nil {
		return nil, err
	}
	return &album, nil
}
