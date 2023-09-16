package nix_search_test

import (
	"context"
	"testing"

	"github.com/luisnquin/nix-search/internal/config"
	nix_search "github.com/luisnquin/nix-search/internal/nix/search"
)

func TestFlakeOptionsSmoke(t *testing.T) {
	appConfig, err := config.Load(true)
	if err != nil {
		t.Fatal(err)
	}

	client := nix_search.NewClient(appConfig)
	ctx := context.Background()

	options, err := client.SearchFlakeOptions(ctx, nix_search.ELASTIC_SEARCH_FLAKES_ID, "wayland", 50)
	if err != nil {
		t.Fatal(err)
	}

	for _, option := range options {
		t.Log(option)
	}
}
