package api

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func (c *Client) GetStationTracks(stationType, tag string, lastTrackID string) (*StationTracksResult, error) {
	params := url.Values{
		"settings2": {"true"},
	}
	if lastTrackID != "" {
		params.Set("queue", lastTrackID)
	}
	path := fmt.Sprintf("/rotor/station/%s:%s/tracks", stationType, tag)
	raw, err := c.get(path, params)
	if err != nil {
		return nil, err
	}
	var result StationTracksResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SendStationFeedback(stationType, tag, batchID, action, trackID string, totalSeconds float64) error {
	path := fmt.Sprintf("/rotor/station/%s:%s/feedback", stationType, tag)
	form := url.Values{
		"type":          {action},
		"timestamp":     {fmt.Sprintf("%.0f", float64(0))},
		"batch-id":      {batchID},
	}
	if trackID != "" {
		form.Set("track-id", trackID)
	}
	if totalSeconds > 0 {
		form.Set("totalPlayedSeconds", fmt.Sprintf("%.0f", totalSeconds))
	}
	_, err := c.post(path, form)
	return err
}
