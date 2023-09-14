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
	switch app.tabs.search.Name {
	case HOME_MANAGER_OPTIONS:
		return app.searchHomeManagerOptions(ctx, input)
	case NIX_PACKAGES:
		return app.searchNixPackages(ctx, input)
	}

	return "", nil
}

func (app App) searchHomeManagerOptions(ctx context.Context, input string) (string, error) {
	options, err := app.nixClient.SearchHomeManagerOptions(ctx, input)
	if err != nil {
		uerr, ok := err.(*url.Error)
		if ok && uerr.Timeout() {
			return "", nil
		}

		return "", err // TODO: send to terminal screen and do not display context cancelled error
	}

	prettyOptions := lo.Map(options, func(opt *nix.HomeManagerOption, _ int) string {
		return fmt.Sprintf("%s - %s\n %s\n Example: %s\n Default: %s\n",
			opt.Title, opt.Description, opt.Position, opt.Example, opt.Default)
	})

	return strings.Join(prettyOptions, "\n\n"), nil
}

func (app App) searchNixOSOptions(ctx context.Context, input string) (string, error) {
	channelStatus := nix.CHANNEL_STATUS_ROLLING

	channel, found := app.config.Internal.Nix.FindChannelWithStatus(channelStatus)
	if !found {
		panic(fmt.Sprintf("unable to find channel with %s status", channelStatus))
	}

	options, err := app.nixClient.SearchNixOSOptions(ctx, channel.Branch, input, 100)
	if err != nil {
		return "", err
	}

	prettyOptions := lo.Map(options, func(option *nix.Option, _ int) string {
		return fmt.Sprintf("%s - %s\nType: %s\nDefault: %s\n", option.Name, option.Description, *option.Example, option.Default)
	})

	return strings.Join(prettyOptions, "\n\n"), nil
}

func (app App) searchNixPackages(ctx context.Context, input string) (string, error) {
	channelStatus := nix.CHANNEL_STATUS_ROLLING

	channel, found := app.config.Internal.Nix.FindChannelWithStatus(channelStatus)
	if !found {
		panic(fmt.Sprintf("unable to find channel with %s status", channelStatus))
	}

	packages, err := app.nixClient.SearchPackages(ctx, channel.Branch, input, 100)
	if err != nil {
		return "", err
	}

	prettyPkgs := lo.Map(packages, func(pkg *nix.Package, _ int) string {
		return fmt.Sprintf("%s (%s) - %s\nPrograms: %v\nOutputs: %v\n%s\n", pkg.Name, pkg.Version,
			pkg.Description, pkg.Outputs, pkg.Programs, app.findSource(channel, *pkg.RepositoryPosition))
	})

	return strings.Join(prettyPkgs, "\n\n"), nil
}

func (app App) findSource(channel config.NixChannel, source string) string {
	return fmt.Sprintf("https://github.com/NixOS/nixpkgs/blob/%s/%s",
		channel.Branch, strings.Replace(source, ":", "#L", -1))
}
