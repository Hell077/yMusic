package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Token  string `json:"token,omitempty"`
	Theme  string `json:"theme,omitempty"`
	Volume int    `json:"volume,omitempty"`
}

var (
	configDir  string
	configPath string
)

func init() {
	home, _ := os.UserHomeDir()
	configDir = filepath.Join(home, ".config", "ymusic")
	configPath = filepath.Join(configDir, "config.json")
}

func ConfigDir() string {
	return configDir
}

func Load() (*Config, error) {
	cfg := &Config{
		Theme:  "dark",
		Volume: 70,
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) Save() error {
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0600)
}
