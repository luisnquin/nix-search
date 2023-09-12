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
		ElasticSearch  ElasticSearchConfig `json:"elastic_search"`
		HomeManager    HomeManagerConfig   `json:"home_manager_options"`
		DefaultChannel string              `json:"default_channel"`
		Channels       []NixChannel        `json:"channels"`
	}

	NixChannel struct {
		ID     string `json:"id"`
		Branch string `json:"branch"`
		JobSet string `json:"jobset"`
		Status string `json:"status"`
	}

	ElasticSearchConfig struct {
		Host           string `json:"host"`
		Username       string `json:"username"`
		Password       string `json:"password"`
		MappingVersion string `json:"mapping_version"`
	}

	HomeManagerConfig struct {
		DataURL string `json:"data_url"`
	}
)

//go:embed internal.config.json
var internalConfig []byte

func Load() *Config {
	var c Config

	if err := json.Unmarshal(internalConfig, &c.Internal); err != nil {
		panic(err)
	}

	return &c
}
