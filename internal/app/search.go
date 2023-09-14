package app

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/luisnquin/nix-search/internal/config"
	"github.com/luisnquin/nix-search/internal/nix"
	"github.com/samber/lo"
)

// Search states.
const (
	MAPPING   = "mapping results"
	SEARCHING = "searching"
	FETCHING  = "fetching"
	WAITING   = "waiting"
)

var ErrChannelNotFound = fmt.Errorf("channel not found")

func (app *App) performSearch(ctx context.Context, input string) {
	app.updateCurrentStatus(SEARCHING)

	results, err := app.performSearchAndGetResults(ctx, input)
	if err != nil { // TODO: handle error and send to widget
		app.widgets.resultsBoard.Reset()

		return
	}

	app.updateCurrentStatus(WAITING)

	app.widgets.resultsBoard.Reset()
	app.widgets.resultsBoard.Write(results)
}

func (app *App) performSearchAndGetResults(ctx context.Context, input string) (string, error) {
	chSearchStatus := make(chan string)

	go func() {
		for status := range chSearchStatus {
			app.updateCurrentStatus(status)
		}
	}()

	switch app.tabs.search.Name {
	case HOME_MANAGER_OPTIONS:
		return app.searchHomeManagerOptions(ctx, input, chSearchStatus)
	case NIX_PACKAGES:
		return app.searchNixPackages(ctx, input, chSearchStatus)
	case NIXOS_OPTIONS:
		return app.searchNixOSOptions(ctx, input, chSearchStatus)
	case FLAKES_OPTIONS:
		return app.searchNixFlakeOptions(ctx, input, chSearchStatus)
	case FLAKES_PACKAGES:
		return app.searchNixFlakePackages(ctx, input, chSearchStatus)
	}

	return "", nil
}

func (app App) searchHomeManagerOptions(ctx context.Context, input string, chStatus chan string) (string, error) {
	defer close(chStatus)

	if app.nixClient.HomeManagerOptionsAlreadyFetched() {
		chStatus <- SEARCHING
	} else {
		chStatus <- FETCHING
	}

	options, err := app.nixClient.SearchHomeManagerOptions(ctx, input)
	if err != nil {
		uerr, ok := err.(*url.Error)
		if ok && uerr.Timeout() {
			return "", nil
		}

		return "", err // TODO: send to terminal screen and do not display context cancelled error
	}

	chStatus <- MAPPING

	prettyOptions := lo.Map(options, func(opt *nix.HomeManagerOption, _ int) string {
		return fmt.Sprintf("%s - %s\n %s\n Example: %s\n Default: %s\n",
			opt.Title, opt.Description, opt.Position, opt.Example, opt.Default)
	})

	chStatus <- WAITING

	return strings.Join(prettyOptions, "\n\n"), nil
}

func (app App) searchNixOSOptions(ctx context.Context, input string, chStatus chan string) (string, error) {
	defer close(chStatus)

	channel, found := app.config.Internal.Nix.FindChannel(app.tabs.search.CurrentChannelID)
	if !found {
		return "", ErrChannelNotFound
	}

	chStatus <- SEARCHING

	options, err := app.nixClient.SearchNixOSOptions(ctx, channel.Branch, input, 100)
	if err != nil {
		return "", err
	}

	chStatus <- MAPPING

	prettyOptions := lo.Map(options, func(option *nix.Option, _ int) string {
		return fmt.Sprintf("%s - %s\nExample: %v\nDefault: %s\n",
			option.Name, option.Description, lo.FromPtrOr(option.Example, "null"), option.Default)
	})

	r := strings.Join(prettyOptions, "\n\n")

	chStatus <- WAITING

	return r, nil
}

func (app App) searchNixPackages(ctx context.Context, input string, chStatus chan string) (string, error) {
	defer close(chStatus)

	chStatus <- SEARCHING

	channel, found := app.config.Internal.Nix.FindChannel(app.tabs.search.CurrentChannelID)
	if !found {
		return "", ErrChannelNotFound
	}

	packages, err := app.nixClient.SearchPackages(ctx, channel.Branch, input, 100)
	if err != nil {
		return "", err
	}

	chStatus <- MAPPING

	prettyPkgs := lo.Map(packages, func(pkg *nix.Package, _ int) string {
		return fmt.Sprintf("%s (%s) - %s\nPrograms: %v\nOutputs: %v\n%s\n", pkg.Name, pkg.Version,
			pkg.Description, pkg.Programs, pkg.Outputs, app.findSource(channel, *pkg.RepositoryPosition))
	})

	r := strings.Join(prettyPkgs, "\n\n")

	chStatus <- WAITING

	return r, nil
}

func (app App) searchNixFlakePackages(ctx context.Context, input string, chStatus chan string) (string, error) {
	defer close(chStatus)

	chStatus <- SEARCHING

	packages, err := app.nixClient.SearchFlakePackages(ctx, app.tabs.search.CurrentChannelID, input, 100)
	if err != nil {
		return "", err
	}

	chStatus <- MAPPING

	prettyPkgs := lo.Map(packages, func(pkg *nix.FlakePackage, _ int) string {
		return fmt.Sprintf("%s (%s) - %s\nFlake: %s\nPrograms: %v\nOutputs: %v\n",
			pkg.Name, pkg.Version, pkg.Description, pkg.Flake.Name, pkg.Programs, pkg.Outputs)
	})

	r := strings.Join(prettyPkgs, "\n\n")

	chStatus <- WAITING

	return r, nil
}

func (app App) searchNixFlakeOptions(ctx context.Context, input string, chStatus chan string) (string, error) {
	defer close(chStatus)

	chStatus <- SEARCHING

	options, err := app.nixClient.SearchFlakeOptions(ctx, app.tabs.search.CurrentChannelID, input, 100)
	if err != nil {
		return "", err
	}

	chStatus <- MAPPING

	prettyOptions := lo.Map(options, func(option *nix.FlakeOption, _ int) string {
		return fmt.Sprintf("%s - %s\nFlake: %s\nExample: %v\nDefault: %s\n",
			option.Name, option.Description, option.Flake.Name, lo.FromPtrOr(option.Example, "null"), option.Default)
	})

	r := strings.Join(prettyOptions, "\n\n")

	chStatus <- WAITING

	return r, nil
}

func (app App) findSource(channel config.NixChannel, source string) string {
	return fmt.Sprintf("https://github.com/NixOS/nixpkgs/blob/%s/%s",
		channel.Branch, strings.Replace(source, ":", "#L", -1))
}
