package nix_search

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/luisnquin/nix-search/internal/nix"
	"github.com/samber/lo"
	"mvdan.cc/xurls/v2"
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

	searchTerm = strings.TrimSpace(searchTerm)

	results := lo.Filter(options, func(option *nix.HomeManagerOption, _ int) bool {
		return strings.HasPrefix(option.Title, searchTerm)
	})
	if len(results) > 0 {
		return results, nil
	}

	return lo.Filter(options, func(option *nix.HomeManagerOption, _ int) bool {
		return strings.Contains(option.Type, searchTerm)
	}), nil
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

	rxStrict := xurls.Strict()
	r := strings.NewReplacer("&lt;", "<", "&gt;", ">")

	for i, option := range data.Options {
		data.Options[i].Position = rxStrict.FindString(option.Position)
		data.Options[i].Title = r.Replace(option.Title)
	}

	return &data, err
}
