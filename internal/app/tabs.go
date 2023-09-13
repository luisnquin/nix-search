package app

import "github.com/samber/lo"

type (
	searchTabConfig struct {
		Tab          searchTab
		Label        string
		Source       string
		Status       string
		Prompt       string
		WaitForEnter bool
	}

	searchTab int
)

// Searcher tabs.
const (
	HOME_MANAGER_OPTIONS searchTab = iota
	FLAKES_PACKAGES
	FLAKES_OPTIONS
	NIXOS_OPTIONS
	NIX_PACKAGES
)

// Tab states.
const (
	SEARCHING = "searching"
	FETCHING  = "fetching"
	WAITING   = "waiting"
)

// Data sources.
const (
	ELASTIC_SEARCH_SOURCE = "elastic search"
	MEMORY_SOURCE         = "memory"
)

const (
	SEARCH_PACKAGES_PROMPT = "Search packages: "
	SEARCH_OPTIONS_PROMPT  = "Search options: "
)

func (a App) getSearchTabs() []searchTabConfig {
	return []searchTabConfig{
		{
			Tab:          HOME_MANAGER_OPTIONS,
			Label:        "Home manager options",
			Source:       MEMORY_SOURCE,
			Status:       WAITING,
			Prompt:       SEARCH_OPTIONS_PROMPT,
			WaitForEnter: false,
		},
		{
			Tab:          NIX_PACKAGES,
			Label:        "Nix packages",
			Source:       ELASTIC_SEARCH_SOURCE,
			Status:       WAITING,
			Prompt:       SEARCH_PACKAGES_PROMPT,
			WaitForEnter: true,
		},
		{
			Tab:          NIXOS_OPTIONS,
			Label:        "NixOS options",
			Source:       ELASTIC_SEARCH_SOURCE,
			Status:       WAITING,
			Prompt:       SEARCH_OPTIONS_PROMPT,
			WaitForEnter: true,
		},
		{
			Tab:          FLAKES_PACKAGES,
			Label:        "Flake packages",
			Source:       ELASTIC_SEARCH_SOURCE,
			Status:       WAITING,
			Prompt:       SEARCH_PACKAGES_PROMPT,
			WaitForEnter: true,
		},
		{
			Tab:          FLAKES_OPTIONS,
			Label:        "Flake options",
			Source:       ELASTIC_SEARCH_SOURCE,
			Status:       WAITING,
			Prompt:       SEARCH_OPTIONS_PROMPT,
			WaitForEnter: true,
		},
	}
}

func (a App) getDefaultSearchTab() *searchTabConfig {
	config, _ := lo.Find(a.getSearchTabs(), func(item searchTabConfig) bool {
		return item.Tab == HOME_MANAGER_OPTIONS
	})

	return &config
}

func (a *App) nextTab() {
	searchTabs := a.getSearchTabs()

	index := a.getCurrentTabIndex()

	if index+1 < len(searchTabs) {
		a.currentSearchTab = &searchTabs[index+1]
	}
}

func (a *App) previousTab() {
	searchTabs := a.getSearchTabs()

	index := a.getCurrentTabIndex()

	if index-1 >= 0 {
		a.currentSearchTab = &searchTabs[index-1]
	}
}

func (a App) getCurrentTabIndex() int {
	_, index, found := lo.FindIndexOf(a.getSearchTabs(), func(searchTab searchTabConfig) bool {
		return searchTab.Tab == a.currentSearchTab.Tab
	})
	if found {
		return index
	}

	return 0
}
