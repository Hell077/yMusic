package api

import "encoding/json"

func (c *Client) GetAccountStatus() (*AccountStatus, error) {
	raw, err := c.get("/account/status", nil)
	if err != nil {
		return nil, err
	}
	var status AccountStatus
	if err := json.Unmarshal(raw, &status); err != nil {
		return nil, err
	}
	return &status, nil
}
