package config

import (
	_ "embed"
	"encoding/json"
)

type Config struct {
	Internal InternalConfig
}

type InternalConfig struct {
	Nix NixConfig `json:"nix"`
}

//go:embed internal.config.json
var internalConfig []byte

func Load() *Config {
	var c Config

	if err := json.Unmarshal(internalConfig, &c.Internal); err != nil {
		panic(err)
	}

	return &c
}
