package app

import (
	"context"
	"net/url"
	"strings"

	"github.com/luisnquin/nix-search/internal/nix"
	"github.com/mum4k/termdash/widgets/text"
	"github.com/samber/lo"
)

func (app App) performSearch(ctx context.Context, input string) {
	results, err := app.performSearchAndGetResults(ctx, input)
	if err != nil { // TODO: handle error and send to widget
		app.resultsBoard.Reset()

		return
	}

	app.resultsBoard.Write(results, text.WriteReplace())
}

func (app App) performSearchAndGetResults(ctx context.Context, input string) (string, error) {
	switch app.currentSearchTab.Name {
	case HOME_MANAGER_OPTIONS:
		return app.searchHomeManagerOptions(ctx, input)
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

	return strings.Join(lo.Map(options, func(opt *nix.HomeManagerOption, _ int) string {
		return opt.String()
	}), " "), nil
}
