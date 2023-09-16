package gui

import (
	"context"
	_ "embed"
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

//go:embed outputs/nix-packages.tpl
var nixPackagesOutputTpl string

var ErrChannelNotFound = fmt.Errorf("channel not found")

func (g *GUI) performSearch(ctx context.Context, input string) {
	g.updateCurrentStatus(SEARCHING)

	results, err := g.performSearchAndGetResults(ctx, input)
	if err != nil { // TODO: handle error and send to widget
		g.widgets.resultsBoard.Reset()

		return
	}

	g.updateCurrentStatus(WAITING)

	g.widgets.resultsBoard.Reset()
	g.widgets.resultsBoard.Write(results)
}

func (g *GUI) performSearchAndGetResults(ctx context.Context, input string) (string, error) {
	defer g.handleProgramPanic()

	statusChan := make(chan string)

	go func() {
		for status := range statusChan {
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

	if g.nixClient.HomeManagerOptionsAlreadyFetched() {
		statusChan <- SEARCHING
	} else {
		statusChan <- FETCHING
	}

	options, err := g.nixClient.SearchHomeManagerOptions(ctx, input)
	if err != nil {
		uerr, ok := err.(*url.Error)
		if ok && uerr.Timeout() {
			return "", nil
		}

		return "", err // TODO: send to terminal screen and do not display context cancelled error
	}

	statusChan <- MAPPING

	prettyOptions := lo.Map(options, func(opt *nix.HomeManagerOption, _ int) string {
		example := lo.If(opt.Example != "", opt.Example).Else("<nothing>")

		noteOrNothing := lo.If(opt.Note != "", fmt.Sprintf("Note: %s\n", opt.Note)).Else("")

		return fmt.Sprintf("%s - %s\n%sType: %s\nExample: %s\nDefault: %s\nSource: %s\n",
			opt.Title, opt.Description, noteOrNothing, opt.Type, example, opt.Default, opt.Position)
	})

	statusChan <- WAITING

	return strings.Join(prettyOptions, "\n\n"), nil
}

func (g GUI) searchNixOSOptions(ctx context.Context, input string, statusChan chan string) (string, error) {
	defer close(statusChan)

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

	prettyOptions := lo.Map(options, func(option *nix.Option, _ int) string {
		return fmt.Sprintf("%s - %s\nExample: %v\nDefault: %s\n",
			option.Name, option.Description, lo.FromPtrOr(option.Example, "null"), option.Default)
	})

	r := strings.Join(prettyOptions, "\n\n")

	statusChan <- WAITING

	return r, nil
}

func (g GUI) searchNixPackages(ctx context.Context, input string, statusChan chan string) (string, error) {
	defer close(statusChan)

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

	text, err := getRenderedText(nix.NIX_PACKAGES, nixPackagesOutputTpl, packages)
	if err != nil {
		return "", err
	}

	statusChan <- WAITING

	return text, nil
}

func (g GUI) searchNixFlakePackages(ctx context.Context, input string, statusChan chan string) (string, error) {
	defer close(statusChan)

	statusChan <- SEARCHING

	packages, err := g.nixClient.SearchFlakePackages(ctx, g.tabs.search.CurrentChannelID, input, 100)
	if err != nil {
		return "", err
	}

	statusChan <- MAPPING

	prettyPkgs := lo.Map(packages, func(pkg *nix.FlakePackage, _ int) string {
		return fmt.Sprintf("%s (%s) - %s\nFlake: %s\nPrograms: %v\nOutputs: %v\n",
			pkg.Name, pkg.Version, pkg.Description, pkg.Flake.Name, pkg.Programs, pkg.Outputs)
	})

	r := strings.Join(prettyPkgs, "\n\n")

	statusChan <- WAITING

	return r, nil
}

func (g GUI) searchNixFlakeOptions(ctx context.Context, input string, statusChan chan string) (string, error) {
	defer close(statusChan)

	statusChan <- SEARCHING

	options, err := g.nixClient.SearchFlakeOptions(ctx, g.tabs.search.CurrentChannelID, input, 100)
	if err != nil {
		return "", err
	}

	statusChan <- MAPPING

	prettyOptions := lo.Map(options, func(option *nix.FlakeOption, _ int) string {
		return fmt.Sprintf("%s - %s\nFlake: %s\nExample: %v\nDefault: %s\n",
			option.Name, option.Description, option.Flake.Name, lo.FromPtrOr(option.Example, "null"), option.Default)
	})

	r := strings.Join(prettyOptions, "\n\n")

	statusChan <- WAITING

	return r, nil
}

func (g GUI) findSource(channel config.NixChannel, source string) string {
	return fmt.Sprintf("https://github.com/NixOS/nixpkgs/blob/%s/%s",
		channel.Branch, strings.Replace(source, ":", "#L", -1))
}
