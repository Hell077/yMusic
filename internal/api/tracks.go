package api

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

const downloadSalt = "XGRlBW9FXlekgbPrRHuSiA"

func (c *Client) GetTrack(id string) (*Track, error) {
	raw, err := c.get("/tracks/"+id, nil)
	if err != nil {
		return nil, err
	}
	var tracks []Track
	if err := json.Unmarshal(raw, &tracks); err != nil {
		return nil, err
	}
	if len(tracks) == 0 {
		return nil, fmt.Errorf("track not found: %s", id)
	}
	return &tracks[0], nil
}

func (c *Client) GetDownloadInfo(trackID string) ([]DownloadInfo, error) {
	raw, err := c.get("/tracks/"+trackID+"/download-info", nil)
	if err != nil {
		return nil, err
	}
	var infos []DownloadInfo
	if err := json.Unmarshal(raw, &infos); err != nil {
		return nil, err
	}
	return infos, nil
}

func (c *Client) GetDirectURL(trackID string) (string, error) {
	infos, err := c.GetDownloadInfo(trackID)
	if err != nil {
		return "", err
	}

	// Pick best quality: prefer mp3 320, then highest bitrate
	var best *DownloadInfo
	for i := range infos {
		info := &infos[i]
		if info.Preview {
			continue
		}
		if best == nil || info.Bitrate > best.Bitrate {
			best = info
		}
	}
	if best == nil {
		return "", fmt.Errorf("no download info for track %s", trackID)
	}

	// Fetch the download data
	srcURL := best.Src
	if !strings.HasPrefix(srcURL, "http") {
		srcURL = "https://" + srcURL
	}
	// Add format=json
	if strings.Contains(srcURL, "?") {
		srcURL += "&format=json"
	} else {
		srcURL += "?format=json"
	}

	data, err := c.getExternal(srcURL)
	if err != nil {
		return "", err
	}

	var dd DownloadData
	if err := json.Unmarshal(data, &dd); err != nil {
		return "", err
	}

	return buildDirectURL(dd), nil
}

func buildDirectURL(dd DownloadData) string {
	// sign = md5("XGRlBW9FXlekgbPrRHuSiA" + path[1:] + s)
	path := dd.Path
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	signData := downloadSalt + path + dd.S
	sign := fmt.Sprintf("%x", md5.Sum([]byte(signData)))

	return fmt.Sprintf("https://%s/get-mp3/%s/%s%s",
		dd.Host,
		sign,
		dd.TS,
		dd.Path,
	)
}

func (c *Client) GetTracks(ids []string) ([]Track, error) {
	form := url.Values{
		"track-ids": {strings.Join(ids, ",")},
	}
	raw, err := c.post("/tracks", form)
	if err != nil {
		return nil, err
	}
	var tracks []Track
	if err := json.Unmarshal(raw, &tracks); err != nil {
		return nil, err
	}
	return tracks, nil
}
