package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/luisnquin/nix-search/internal"
)

type Config struct {
	Internal InternalConfig
	LogsPath string
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

	if c.LogsPath == "" {
		c.LogsPath = filepath.Join(
			os.TempDir(), fmt.Sprintf("%s.log", internal.PROGRAM_NAME),
		)
	}

	return &c
}
