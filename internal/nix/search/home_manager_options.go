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
	searchTerm = strings.TrimSpace(searchTerm)

	options, err := c.getHomeManagerOptions(ctx)
	if err != nil {
		return nil, err
	}

	matchFns := []func(s string, otherS string) bool{
		strings.HasPrefix, strings.Contains, strings.HasSuffix,
	}

	for _, matchFn := range matchFns {
		results := lo.Filter(options, func(option *nix.HomeManagerOption, _ int) bool {
			return matchFn(option.Title, searchTerm)
		})
		if len(results) > 0 {
			return results, nil
		}
	}

	return nil, nil
}

func (c *Client) HomeManagerOptionsAlreadyFetched() bool {
	return c.store.homeManagerShell.data != nil
}

func (c *Client) getHomeManagerOptions(ctx context.Context) ([]*nix.HomeManagerOption, error) {
	var err error

	c.store.homeManagerShell.Once.Do(func() {
		c.store.homeManagerShell.data, err = c.fetchHomeManagerOptions(ctx)
	})
	if err != nil {
		return nil, err
	}

	return c.store.homeManagerShell.data.Options, nil
}

func (c Client) fetchHomeManagerOptions(ctx context.Context) (*homeManagerOptionsData, error) {
	optionsUrl := c.config.Internal.Nix.Sources.HomeManagerOptions.URL

	response, err := doGET(ctx, optionsUrl)
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