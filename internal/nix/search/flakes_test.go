package nix_search_test

import (
	"context"
	"testing"

	"github.com/luisnquin/nix-search/internal/config"
	nix_search "github.com/luisnquin/nix-search/internal/nix/search"
)

func TestFlakeOptionsSmoke(t *testing.T) {
	ctx, appConfig := context.Background(), config.Load()

	client := nix_search.NewClient(appConfig)

	options, err := client.SearchFlakeOptions(ctx, "wayland", 50)
	if err != nil {
		t.Fatal(err)
	}

	for _, option := range options {
		t.Log(option)
	}
}

func TestFlakePackagesSmoke(t *testing.T) {
	ctx, appConfig := context.Background(), config.Load()

	client := nix_search.NewClient(appConfig)

	options, err := client.SearchFlakePackages(ctx, "wayland", 50)
	if err != nil {
		t.Fatal(err)
	}

	for _, option := range options {
		t.Log(option)
	}
}
