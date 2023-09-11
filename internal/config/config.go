package config

import (
	_ "embed"
	"encoding/json"
)

type Config struct {
	Internal InternalConfig
}

type (
	InternalConfig struct {
		NixOSElasticSearch NixOSElasticSearchConfig `json:"nixos_elastic_search"`
		HomeManagerOptions HomeManagerOptionsConfig `json:"home_manager_options"`
	}

	NixOSElasticSearchConfig struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	HomeManagerOptionsConfig struct {
		DataURL string `json:"data_url"`
	}
)

//go:embed internal-config.json
var internalConfig []byte

func Load() *Config {
	var c Config

	if err := json.Unmarshal(internalConfig, &c.Internal); err != nil {
		panic(err)
	}

	return &c
}
