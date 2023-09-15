package gui

import (
	"github.com/luisnquin/nix-search/internal/config"
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

func (g *GUI) getSearchTabs() []searchTabConfig {
	channelIds := lo.Map(g.config.Internal.Nix.Channels, func(config config.NixChannel, _ int) string {
		return config.ID
	})

	return []searchTabConfig{
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
}

func (g GUI) getDefaultSearchTab() *searchTabConfig {
	config, _ := lo.Find(g.getSearchTabs(), func(item searchTabConfig) bool {
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
		g.tabs.search = &searchTabs[index+1]
		g.updateWidgetTexts()
	}
}

func (g *GUI) previousTab() {
	searchTabs := g.getSearchTabs()
	index := g.getCurrentTabIndex()

	if index-1 >= 0 {
		g.tabs.search = &searchTabs[index-1]
		g.updateWidgetTexts()
	}
}
