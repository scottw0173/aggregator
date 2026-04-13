package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DBURL    string `json:"db_url"`
	UserName string `json:"user_name"`
}

func Read() (Config, error) {
	configPath, err := getCfgFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (cfg *Config) SetUser(name string) error {
	cfg.UserName = name
	configPath, err := getCfgFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getCfgFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	config_string := ".gatorconfig.json"
	configPath := filepath.Join(home, config_string)
	return configPath, nil
}
