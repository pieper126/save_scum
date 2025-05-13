package internal

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ReadFrom string `json:"readFrom"`
	SaveTo   string `json:"saveTo"`
}

func (config *Config) Validate() error {
	if _, err := os.Stat(config.ReadFrom); os.IsNotExist(err) {
		return fmt.Errorf("readFrom does not exist: %s", err)
	}

	if _, err := os.Stat(config.SaveTo); os.IsNotExist(err) {
		return fmt.Errorf("saveTo does not exist: %s", err)
	}

	return nil
}

func LoadConfig(filename string) (Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, fmt.Errorf("error loading config: %w", err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, fmt.Errorf("something wrong with the config file: %w", err)
	}

	return config, nil
}
