package nix_search

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/luisnquin/nix-search/internal/nix"
)

type homeManagerOptionsData struct {
	LastUpdate string                   `json:"last_update"`
	Options    []*nix.HomeManagerOption `json:"options"`
}

func (c *Client) SearchHomeManagerOptions(ctx context.Context, searchTerm string) ([]*nix.HomeManagerOption, error) {
	options, err := c.getHomeManagerOptions(ctx)
	if err != nil {
		return nil, err
	}

	var results []*nix.HomeManagerOption

	for _, option := range options {
		if strings.HasPrefix(option.Title, searchTerm) {
			results = append(results, option)
		}
	}

	return results, nil
}

func (c *Client) getHomeManagerOptions(ctx context.Context) ([]*nix.HomeManagerOption, error) {
	var err error

	c.store.homeManagerShell.Once.Do(func() {
		c.store.homeManagerShell.data, err = c.fetchHomeManagerOptions()
	})
	if err != nil {
		return nil, err
	}

	return c.store.homeManagerShell.data.Options, nil
}

func (c Client) fetchHomeManagerOptions() (*homeManagerOptionsData, error) {
	optionsUrl := c.config.Internal.Nix.Sources.HomeManagerOptions.URL

	response, err := doGET(optionsUrl)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var data homeManagerOptionsData

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, err
}
