package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ProtonMail/go-appdir"
	"github.com/luisnquin/nix-search/internal"
	"github.com/luisnquin/nix-search/internal/nix"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Internal InternalConfig
	// The path to the logs file of the program.
	LogFile    string `yaml:"log_file"`
	SearchTabs struct {
		// The desired order of the search tabs.
		Order []string `yaml:"order"`
		// The name of the default tab.
		Selected string `yaml:"selected"`
	} `yaml:"search_tabs"`
}

type InternalConfig struct {
	Nix NixConfig `json:"nix"`
}

//go:embed internal.config.json
var internalConfig []byte

func Load(test bool) (*Config, error) {
	var c Config

	if !test {
		if err := c.loadUserConfig(); err != nil {
			return nil, err
		}

		if err := c.validateUserConfig(); err != nil {
			return nil, err
		}
	}

	if err := json.Unmarshal(internalConfig, &c.Internal); err != nil {
		return nil, err
	}

	if c.LogFile == "" {
		c.LogFile = filepath.Join(
			os.TempDir(), fmt.Sprintf("%s.log", internal.PROGRAM_NAME),
		)
	}

	return &c, nil
}

func (c *Config) loadUserConfig() error {
	dirs := appdir.New(internal.PROGRAM_NAME)

	for _, fileName := range []string{"config.yaml", "config.yml"} {
		filePath := filepath.Join(dirs.UserConfig(), fileName)

		info, err := os.Stat(filePath)
		if err == nil && !info.IsDir() {
			content, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}

			return yaml.Unmarshal(content, c)
		}
	}

	return nil
}

func (c *Config) validateUserConfig() error {
	if c.SearchTabs.Selected != "" && !lo.Contains(nix.GetSourceNames(), c.SearchTabs.Selected) {
		return fmt.Errorf("unknown tab: %s", c.SearchTabs.Selected)
	}

	tabsOrder := lo.Uniq(c.SearchTabs.Order)

	unknownTabs := lo.Filter(tabsOrder, func(name string, index int) bool {
		return !lo.Contains(nix.GetSourceNames(), name)
	})
	if l := len(unknownTabs); l > 0 {
		if l == 1 {
			return fmt.Errorf("unknown tab: %s", unknownTabs[0])
		}

		return fmt.Errorf("unknown tabs: %s", strings.Join(unknownTabs, ", "))
	}

	c.SearchTabs.Order = tabsOrder

	return nil
}
