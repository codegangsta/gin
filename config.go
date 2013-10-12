package gin

import (
	"encoding/json"
	"fmt"
	"os"
)

type ServerConfig struct {
	Port int `json:"port"`
}

type Config struct {
	Port   int          `json:"port"`
	Server ServerConfig `json:"server"`
}

func LoadConfig(path string) (*Config, error) {
	configFile, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("Unable to read configuration file %s\n%s", path, err.Error())
	}

	config := new(Config)

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse configuration file %s\n%s", path, err.Error())
	}

	return config, nil
}
