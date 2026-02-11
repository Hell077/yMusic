package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func (c *Client) GetLikedTracks(uid int) (*LikesResult, error) {
	raw, err := c.get(fmt.Sprintf("/users/%d/likes/tracks", uid), nil)
	if err != nil {
		return nil, err
	}
	var result LikesResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type LikedAlbum struct {
	ID        int    `json:"id"`
	Timestamp string `json:"timestamp"`
}

func (c *Client) GetLikedAlbums(uid int) ([]LikedAlbum, error) {
	raw, err := c.get(fmt.Sprintf("/users/%d/likes/albums", uid), nil)
	if err != nil {
		return nil, err
	}
	var result []LikedAlbum
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetAlbums(ids []string) ([]Album, error) {
	form := url.Values{
		"album-ids": {strings.Join(ids, ",")},
	}
	raw, err := c.post("/albums", form)
	if err != nil {
		return nil, err
	}
	var albums []Album
	if err := json.Unmarshal(raw, &albums); err != nil {
		return nil, err
	}
	return albums, nil
}

func (c *Client) LikeTrack(uid int, trackID string) error {
	form := url.Values{"track-ids": {trackID}}
	_, err := c.post(fmt.Sprintf("/users/%d/likes/tracks/add-multiple", uid), form)
	return err
}

func (c *Client) UnlikeTrack(uid int, trackID string) error {
	form := url.Values{"track-ids": {trackID}}
	_, err := c.post(fmt.Sprintf("/users/%d/likes/tracks/remove", uid), form)
	return err
}
