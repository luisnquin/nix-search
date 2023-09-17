package gui

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/luisnquin/nix-search/internal/nix"
	nix_search "github.com/luisnquin/nix-search/internal/nix/search"
)

// Search states.
const (
	MAPPING   = "mapping results"
	SEARCHING = "searching"
	FETCHING  = "fetching"
	WAITING   = "waiting"
)

var (
	//go:embed outputs/nix-packages.tpl
	nixPackagesOutputTpl string
	//go:embed outputs/nixos-options.tpl
	nioxsOptionsOutputTpl string
	//go:embed outputs/home-options.tpl
	homeOptionsOutputTpl string
	//go:embed outputs/flake-options.tpl
	flakeOptionsOutputTpl string
	//go:embed outputs/flake-packages.tpl
	flakePackagesOutputTpl string
)

var ErrChannelNotFound = fmt.Errorf("channel not found")

func (g *GUI) performSearch(ctx context.Context, input string) {
	results, err := g.performSearchAndGetResults(ctx, input)
	if err != nil {
		g.widgets.resultsBoard.Reset()
		g.logger.Err(err).Msg("search failed with an error...")
	} else {
		g.widgets.resultsBoard.Reset()
		g.logger.Debug().Str("results", results).Send()

		if len(results) == 0 {
			g.widgets.resultsBoard.Write("<no results>")
		} else {
			g.widgets.resultsBoard.Write(results)
		}
	}
}

func (g *GUI) performSearchAndGetResults(ctx context.Context, input string) (string, error) {
	defer g.handleProgramPanic()

	statusChan := make(chan string)

	go func() {
		for status := range statusChan {
			g.logger.Trace().Str("search-status", status)
			g.updateCurrentStatus(status)
		}
	}()

	switch g.tabs.search.Name {
	case HOME_MANAGER_OPTIONS:
		return g.searchHomeManagerOptions(ctx, input, statusChan)
	case NIX_PACKAGES:
		return g.searchNixPackages(ctx, input, statusChan)
	case NIXOS_OPTIONS:
		return g.searchNixOSOptions(ctx, input, statusChan)
	case FLAKES_OPTIONS:
		return g.searchNixFlakeOptions(ctx, input, statusChan)
	case FLAKES_PACKAGES:
		return g.searchNixFlakePackages(ctx, input, statusChan)
	}

	return "", nil
}

func (g GUI) searchHomeManagerOptions(ctx context.Context, input string, statusChan chan string) (string, error) {
	defer close(statusChan)
	defer g.updateCurrentStatus(WAITING)

	if g.nixClient.HomeManagerOptionsAlreadyFetched() {
		statusChan <- SEARCHING
	} else {
		statusChan <- FETCHING
	}

	options, err := g.nixClient.SearchHomeManagerOptions(ctx, input)
	if err != nil {
		return "", err // TODO: send to terminal screen and do not display context cancelled error
	}

	statusChan <- MAPPING

	return getRenderedText(nix.HOME_OPTIONS, "", homeOptionsOutputTpl, options)
}

func (g GUI) searchNixOSOptions(ctx context.Context, input string, statusChan chan string) (string, error) {
	defer close(statusChan)
	defer g.updateCurrentStatus(WAITING)

	channel, found := g.config.Internal.Nix.FindChannel(g.tabs.search.CurrentChannelID)
	if !found {
		return "", ErrChannelNotFound
	}

	statusChan <- SEARCHING

	options, err := g.nixClient.SearchNixOSOptions(ctx, channel.Branch, input, 100)
	if err != nil {
		return "", err
	}

	statusChan <- MAPPING

	return getRenderedText(nix.NIXOS_OPTIONS, channel.Branch, nioxsOptionsOutputTpl, options)
}

func (g GUI) searchNixPackages(ctx context.Context, input string, statusChan chan string) (string, error) {
	defer close(statusChan)
	defer g.updateCurrentStatus(WAITING)

	statusChan <- SEARCHING

	channel, found := g.config.Internal.Nix.FindChannel(g.tabs.search.CurrentChannelID)
	if !found {
		return "", ErrChannelNotFound
	}

	packages, err := g.nixClient.SearchPackages(ctx, channel.Branch, input, 100)
	if err != nil {
		return "", err
	}

	statusChan <- MAPPING

	return getRenderedText(nix.NIX_PACKAGES, channel.Branch, nixPackagesOutputTpl, packages)
}

func (g GUI) searchNixFlakePackages(ctx context.Context, input string, statusChan chan string) (string, error) {
	defer close(statusChan)
	defer g.updateCurrentStatus(WAITING)

	statusChan <- SEARCHING

	packages, err := g.nixClient.SearchFlakePackages(ctx, g.tabs.search.CurrentChannelID, input, 100)
	if err != nil {
		return "", err
	}

	statusChan <- MAPPING

	return getRenderedText(nix.FLAKE_PACKAGES, nix_search.ELASTIC_SEARCH_FLAKES_ID, flakePackagesOutputTpl, packages)
}

func (g GUI) searchNixFlakeOptions(ctx context.Context, input string, statusChan chan string) (string, error) {
	defer close(statusChan)
	defer g.updateCurrentStatus(WAITING)

	statusChan <- SEARCHING

	options, err := g.nixClient.SearchFlakeOptions(ctx, g.tabs.search.CurrentChannelID, input, 100)
	if err != nil {
		return "", err
	}

	statusChan <- MAPPING

	return getRenderedText(nix.FLAKE_OPTIONS, "", flakeOptionsOutputTpl, options)
}
