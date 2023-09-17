package nix_search

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/luisnquin/nix-search/internal/nix"
	"github.com/samber/lo"
	"mvdan.cc/xurls/v2"
)

type (
	homeManagerOptionsData struct {
		LastUpdate string              `json:"last_update"`
		Options    []homeManagerOption `json:"options"`
	}

	// The DTO used to map the home manager response items.
	homeManagerOption struct {
		// The option name (e.g. programs.zsh.enable)
		Title string `json:"title"`
		// The option description.
		Description string `json:"description"`
		// An additional maintainer note about the option.
		Note string `json:"note"`
		// The expected(s) data type(s) of the option.
		Type string `json:"type"`
		// The default value of the option if allowed.
		Default string `json:"default"`
		// The example value of the option.
		Example string `json:"example"`
		// The repository file where the option was declared.
		Position string `json:"declared_by"`
	}
)

func (c *Client) SearchHomeManagerOptions(ctx context.Context, searchTerm string) ([]*nix.Option, error) {
	searchTerm = strings.TrimSpace(searchTerm)

	options, err := c.getHomeManagerOptions(ctx)
	if err != nil {
		return nil, err
	}

	matchFns := []func(s string, otherS string) bool{
		strings.HasPrefix, strings.Contains, strings.HasSuffix,
	}

	for i := range matchFns {
		matchFn := matchFns[i]

		results := lo.Filter(options, func(option *nix.Option, _ int) bool {
			return matchFn(option.Name, searchTerm)
		})
		if len(results) > 0 {
			return results, nil
		}
	}

	return nil, nil
}

func (c *Client) HomeManagerOptionsAlreadyFetched() bool {
	return c.store.homeManagerShell.options != nil
}

func (c *Client) getHomeManagerOptions(ctx context.Context) ([]*nix.Option, error) {
	var err error

	c.store.homeManagerShell.Once.Do(func() {
		c.store.homeManagerShell.options, err = c.fetchHomeManagerOptions(ctx)
	})

	if err != nil {
		return nil, err
	}

	return c.store.homeManagerShell.options, nil
}

func (c Client) fetchHomeManagerOptions(ctx context.Context) ([]*nix.Option, error) {
	optionsUrl := c.config.Internal.Nix.Sources.HomeManagerOptions.URL

	response, err := doGET(ctx, optionsUrl)
	if err != nil {
		uerr, ok := err.(*url.Error)
		if ok && uerr.Timeout() {
			return nil, nil
		}

		return nil, err
	}

	var data homeManagerOptionsData

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		response.Body.Close()

		return nil, err
	}

	if err := response.Body.Close(); err != nil {
		return nil, err
	}

	rxStrict := xurls.Strict()
	r := strings.NewReplacer("&lt;", "<", "&gt;", ">")

	options := make([]*nix.Option, len(data.Options))

	for i, option := range data.Options {
		example := strings.TrimFunc(option.Example, func(r rune) bool {
			return r == '\n'
		})

		options[i] = &nix.Option{
			Name:            r.Replace(option.Title),
			Description:     option.Description,
			LongDescription: option.Note,
			Type:            option.Type,
			Default:         option.Default,
			Example:         &example,
			Source:          lo.ToPtr(rxStrict.FindString(option.Position)),
		}
	}

	return options, err
}
