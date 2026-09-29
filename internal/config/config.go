package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	gatorconfigPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	content, err := os.ReadFile(gatorconfigPath)
	if err != nil {
		return Config{}, err
	}

	gatorconfig := Config{}
	err = json.Unmarshal(content, &gatorconfig)
	if err != nil {
		return Config{}, err
	}

	return gatorconfig, nil
}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	return write(*cfg)
}

func write(cfg Config) error {
	gatorconfigPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	cfgJson, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(gatorconfigPath, cfgJson, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(home, configFileName)
	return fullPath, nil
}
