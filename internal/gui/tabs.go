package gui

import (
	"sort"

	"github.com/luisnquin/nix-search/internal/config"
	"github.com/luisnquin/nix-search/internal/nix"
	nix_search "github.com/luisnquin/nix-search/internal/nix/search"
	"github.com/samber/lo"
)

type searchTabConfig struct {
	Name             searchTab
	Label            string
	Source           string
	WaitForEnter     bool
	ChannelIDs       []string
	CurrentChannelID string
}

type searchTab int

// Searcher tabs.
const (
	HOME_MANAGER_OPTIONS searchTab = iota
	FLAKES_PACKAGES
	FLAKES_OPTIONS
	NIXOS_OPTIONS
	NIX_PACKAGES
)

// Data sources.
const (
	ELASTIC_SEARCH_SOURCE = "Elastic Search"
	MEMORY_SOURCE         = "External file(in-memory)"
)

func (t searchTab) String() string {
	switch t {
	case HOME_MANAGER_OPTIONS:
		return nix.HOME_OPTIONS
	case FLAKES_PACKAGES:
		return nix.FLAKE_PACKAGES
	case FLAKES_OPTIONS:
		return nix.FLAKE_OPTIONS
	case NIXOS_OPTIONS:
		return nix.NIXOS_OPTIONS
	case NIX_PACKAGES:
		return nix.NIX_PACKAGES
	default:
		return "unknown"
	}
}

func (g *GUI) getSearchTabs() []searchTabConfig {
	channelIds := lo.Map(g.config.Internal.Nix.Channels, func(config config.NixChannel, _ int) string {
		return config.ID
	})

	tabs := []searchTabConfig{
		{
			Name:             NIX_PACKAGES,
			Label:            "Nix packages",
			Source:           ELASTIC_SEARCH_SOURCE,
			WaitForEnter:     true,
			ChannelIDs:       channelIds,
			CurrentChannelID: g.config.Internal.Nix.DefaultChannel,
		},
		{
			Name:         HOME_MANAGER_OPTIONS,
			Label:        "Home manager options",
			Source:       MEMORY_SOURCE,
			WaitForEnter: false,
			ChannelIDs:   nil,
		},
		{
			Name:             NIXOS_OPTIONS,
			Label:            "NixOS options",
			Source:           ELASTIC_SEARCH_SOURCE,
			WaitForEnter:     true,
			ChannelIDs:       channelIds,
			CurrentChannelID: g.config.Internal.Nix.DefaultChannel,
		},
		{
			Name:             FLAKES_PACKAGES,
			Label:            "Flake packages",
			Source:           ELASTIC_SEARCH_SOURCE,
			WaitForEnter:     true,
			ChannelIDs:       []string{nix_search.ELASTIC_SEARCH_FLAKES_ID},
			CurrentChannelID: nix_search.ELASTIC_SEARCH_FLAKES_ID,
		},
		{
			Name:             FLAKES_OPTIONS,
			Label:            "Flake options",
			Source:           ELASTIC_SEARCH_SOURCE,
			WaitForEnter:     true,
			ChannelIDs:       []string{nix_search.ELASTIC_SEARCH_FLAKES_ID},
			CurrentChannelID: nix_search.ELASTIC_SEARCH_FLAKES_ID,
		},
	}

	searchTabsOrder := g.config.SearchTabs.Order

	sort.SliceStable(tabs, func(i, j int) bool {
		iIndex := lo.IndexOf(searchTabsOrder, tabs[i].Name.String())
		jIndex := lo.IndexOf(searchTabsOrder, tabs[j].Name.String())

		return iIndex < jIndex
	})

	return tabs
}

func (g GUI) getSelectedOrDefaultTab() *searchTabConfig {
	searchTabs := g.getSearchTabs()

	if g.config.SearchTabs.Selected != "" {
		config, found := lo.Find(searchTabs, func(config searchTabConfig) bool {
			return config.Name.String() == g.config.SearchTabs.Selected
		})
		if found {
			return &config
		}
	}

	config, _ := lo.Find(searchTabs, func(item searchTabConfig) bool {
		return item.Name == NIX_PACKAGES
	})

	return &config
}

func (g GUI) getCurrentTabIndex() int {
	_, index, found := lo.FindIndexOf(g.getSearchTabs(), func(searchTab searchTabConfig) bool {
		return searchTab.Name == g.tabs.search.Name
	})
	if found {
		return index
	}

	return 0
}

func (g *GUI) nextTab() {
	searchTabs := g.getSearchTabs()
	index := g.getCurrentTabIndex()

	if index+1 < len(searchTabs) {
		tab := searchTabs[index+1]

		g.logger.Trace().Msgf("search tab '%s' -> '%s'", g.tabs.search.Name, tab.Name)

		g.tabs.search = &tab
		g.updateWidgetTexts()
	}
}

func (g *GUI) previousTab() {
	searchTabs := g.getSearchTabs()
	index := g.getCurrentTabIndex()

	if index-1 >= 0 {
		tab := searchTabs[index-1]

		g.logger.Trace().Msgf("search tab '%s' -> '%s'", g.tabs.search.Name, tab.Name)

		g.tabs.search = &tab
		g.updateWidgetTexts()
	}
}
