package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var debugLog *log.Logger

func init() {
	f, err := os.OpenFile("/tmp/ymusic.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		debugLog = log.New(io.Discard, "", 0)
	} else {
		debugLog = log.New(f, "", log.LstdFlags)
	}
}

const baseURL = "https://api.music.yandex.net"

type Client struct {
	token  string
	http   *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		token: token,
		http: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) get(path string, params url.Values) (json.RawMessage, error) {
	u := baseURL + path
	if params != nil {
		u += "?" + params.Encode()
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) post(path string, form url.Values) (json.RawMessage, error) {
	var body io.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	}
	req, err := http.NewRequest("POST", baseURL+path, body)
	if err != nil {
		return nil, err
	}
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return c.do(req)
}

func (c *Client) do(req *http.Request) (json.RawMessage, error) {
	req.Header.Set("Authorization", "OAuth "+c.token)
	req.Header.Set("X-Yandex-Music-Client", "YandexMusicAndroid/24023621")
	req.Header.Set("User-Agent", "Yandex-Music-API")

	debugLog.Printf("API request: %s %s", req.Method, req.URL)

	resp, err := c.http.Do(req)
	if err != nil {
		debugLog.Printf("API error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	debugLog.Printf("API response: %s %d (%d bytes)", req.URL.Path, resp.StatusCode, len(data))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		debugLog.Printf("API error body: %s", string(data[:min(len(data), 500)]))
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(data))
	}

	var envelope struct {
		Result json.RawMessage `json:"result"`
		Error  *struct {
			Name    string `json:"name"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	if envelope.Error != nil {
		return nil, fmt.Errorf("API error: %s: %s", envelope.Error.Name, envelope.Error.Message)
	}

	return envelope.Result, nil
}

func (c *Client) getExternal(rawURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
