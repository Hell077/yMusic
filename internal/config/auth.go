package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	musicClientID     = "23cabbbdc6cd418abb4b39c32c41195d"
	musicClientSecret = "53bc75238f0c4d08a118e51fe9203300"
	tokenURL          = "https://oauth.yandex.ru/token"

	// Authorization code flow — code appears as ?code=XXX in the URL (visible, persistent)
	AuthURL = "https://oauth.yandex.ru/authorize?response_type=code&client_id=" + musicClientID
)

// ExchangeCode exchanges an authorization code for an access token.
func ExchangeCode(code string) (string, error) {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {musicClientID},
		"client_secret": {musicClientSecret},
	}

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
		ErrorDesc   string `json:"error_description"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode token response: %w", err)
	}

	if result.Error != "" {
		return "", fmt.Errorf("%s: %s", result.Error, result.ErrorDesc)
	}
	if result.AccessToken == "" {
		return "", fmt.Errorf("empty access_token in response")
	}

	return result.AccessToken, nil
}
