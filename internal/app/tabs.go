package app

type (
	searchTabConfig struct {
		Tab    searchTab
		Label  string
		Source string
		State  string
		Prompt string
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
	ELASTIC_SEARCH = "elastic search"
	CACHED_FILE    = "cached file"
)

const (
	SEARCH_PACKAGES_PROMPT = "Search packages: "
	SEARCH_OPTIONS_PROMPT  = "Search options: "
)

func (a App) getSearchTabs() []searchTabConfig {
	return []searchTabConfig{
		{
			Tab:    HOME_MANAGER_OPTIONS,
			Label:  "Home manager options",
			Source: CACHED_FILE,
			State:  WAITING,
			Prompt: SEARCH_OPTIONS_PROMPT,
		},
		{
			Tab:    NIX_PACKAGES,
			Label:  "Nix packages",
			Source: ELASTIC_SEARCH,
			State:  WAITING,
			Prompt: SEARCH_PACKAGES_PROMPT,
		},
		{
			Tab:    NIXOS_OPTIONS,
			Label:  "NixOS options",
			Source: ELASTIC_SEARCH,
			State:  WAITING,
			Prompt: SEARCH_OPTIONS_PROMPT,
		},
		{
			Tab:    FLAKES_PACKAGES,
			Label:  "Flake packages",
			Source: ELASTIC_SEARCH,
			State:  WAITING,
			Prompt: SEARCH_PACKAGES_PROMPT,
		},
		{
			Tab:    FLAKES_OPTIONS,
			Label:  "Flake options",
			Source: ELASTIC_SEARCH,
			State:  WAITING,
			Prompt: SEARCH_OPTIONS_PROMPT,
		},
	}
}
