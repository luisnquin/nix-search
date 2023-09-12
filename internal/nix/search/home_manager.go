package nix_search

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/luisnquin/nix-search/internal/nix"
)

type homeManagerOptionsData struct {
	LastUpdate string                   `json:"last_update"`
	Options    []*nix.HomeManagerOption `json:"options"`
}

func (c *Client) SearchHomeManagerOptions(searchTerm string) []*nix.HomeManagerOption {
	c.homeManager.Do(c.tryToFetchHomeManagerOptions)

	var results []*nix.HomeManagerOption

	for _, option := range c.homeManager.data.Options {
		if strings.HasPrefix(option.Title, searchTerm) {
			results = append(results, option)
		}
	}

	return results
}

func (c *Client) tryToFetchHomeManagerOptions() {
	c.homeManager.mu.Lock()
	defer c.homeManager.mu.Unlock()

	data, err := c.fetchHomeManagerOptions()
	if err != nil {
		panic(err)
	}

	c.homeManager.data = data
}

func (c Client) fetchHomeManagerOptions() (homeManagerOptionsData, error) {
	httpClient := http.Client{Timeout: CLIENT_TIMEOUT}
	optionsUrl := c.config.Internal.Nix.Sources.HomeManagerOptions.URL

	r, err := http.NewRequest(http.MethodGet, optionsUrl, http.NoBody)
	if err != nil {
		return homeManagerOptionsData{}, err
	}

	defer r.Body.Close()

	response, err := httpClient.Do(r)
	if err != nil {
		return homeManagerOptionsData{}, err
	}

	defer response.Body.Close()

	var data homeManagerOptionsData

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return homeManagerOptionsData{}, err
	}

	return data, err
}
