package api

import (
	"encoding/json"
	"fmt"
)

func (c *Client) GetUserPlaylists(uid int) ([]Playlist, error) {
	raw, err := c.get(fmt.Sprintf("/users/%d/playlists/list", uid), nil)
	if err != nil {
		return nil, err
	}
	var playlists []Playlist
	if err := json.Unmarshal(raw, &playlists); err != nil {
		return nil, err
	}
	return playlists, nil
}

func (c *Client) GetPlaylist(uid, kind int) (*Playlist, error) {
	raw, err := c.get(fmt.Sprintf("/users/%d/playlists/%d", uid, kind), nil)
	if err != nil {
		return nil, err
	}
	var playlist Playlist
	if err := json.Unmarshal(raw, &playlist); err != nil {
		return nil, err
	}
	return &playlist, nil
}
