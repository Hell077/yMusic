package api

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func (c *Client) Search(query string, page int) (*SearchResult, error) {
	params := url.Values{
		"text":    {query},
		"nocorrect": {"false"},
		"type":    {"all"},
		"page":    {fmt.Sprintf("%d", page)},
	}
	raw, err := c.get("/search", params)
	if err != nil {
		return nil, err
	}
	var result SearchResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SearchTracks(result *SearchResult) ([]Track, error) {
	if result.Tracks == nil || len(result.Tracks.Results) == 0 {
		return nil, nil
	}
	var tracks []Track
	if err := json.Unmarshal(result.Tracks.Results, &tracks); err != nil {
		return nil, err
	}
	return tracks, nil
}

func (c *Client) SearchAlbums(result *SearchResult) ([]Album, error) {
	if result.Albums == nil || len(result.Albums.Results) == 0 {
		return nil, nil
	}
	var albums []Album
	if err := json.Unmarshal(result.Albums.Results, &albums); err != nil {
		return nil, err
	}
	return albums, nil
}

func (c *Client) SearchArtists(result *SearchResult) ([]Artist, error) {
	if result.Artists == nil || len(result.Artists.Results) == 0 {
		return nil, nil
	}
	var artists []Artist
	if err := json.Unmarshal(result.Artists.Results, &artists); err != nil {
		return nil, err
	}
	return artists, nil
}
