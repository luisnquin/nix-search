package app

import "github.com/samber/lo"

type searchTabConfig struct {
	Name         searchTab
	Label        string
	Source       string
	Status       string
	WaitForEnter bool
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

func (a *App) getSearchTabs() []searchTabConfig {
	return []searchTabConfig{
		{
			Name:         NIX_PACKAGES,
			Label:        "Nix packages",
			Source:       ELASTIC_SEARCH_SOURCE,
			Status:       WAITING,
			WaitForEnter: true,
		},
		{
			Name:         HOME_MANAGER_OPTIONS,
			Label:        "Home manager options",
			Source:       MEMORY_SOURCE,
			Status:       WAITING,
			WaitForEnter: false,
		},
		{
			Name:         NIXOS_OPTIONS,
			Label:        "NixOS options",
			Source:       ELASTIC_SEARCH_SOURCE,
			Status:       WAITING,
			WaitForEnter: true,
		},
		{
			Name:         FLAKES_PACKAGES,
			Label:        "Flake packages",
			Source:       ELASTIC_SEARCH_SOURCE,
			Status:       WAITING,
			WaitForEnter: true,
		},
		{
			Name:         FLAKES_OPTIONS,
			Label:        "Flake options",
			Source:       ELASTIC_SEARCH_SOURCE,
			Status:       WAITING,
			WaitForEnter: true,
		},
	}
}

func (a App) getDefaultSearchTab() *searchTabConfig {
	config, _ := lo.Find(a.getSearchTabs(), func(item searchTabConfig) bool {
		return item.Name == NIX_PACKAGES
	})

	return &config
}

func (app App) getCurrentTabIndex() int {
	_, index, found := lo.FindIndexOf(app.getSearchTabs(), func(searchTab searchTabConfig) bool {
		return searchTab.Name == app.tabs.search.Name
	})
	if found {
		return index
	}

	return 0
}

func (app *App) nextTab() {
	searchTabs := app.getSearchTabs()
	index := app.getCurrentTabIndex()

	if index+1 < len(searchTabs) {
		app.tabs.search = &searchTabs[index+1]
	}
}

func (app *App) previousTab() {
	searchTabs := app.getSearchTabs()
	index := app.getCurrentTabIndex()

	if index-1 >= 0 {
		app.tabs.search = &searchTabs[index-1]
	}
}
