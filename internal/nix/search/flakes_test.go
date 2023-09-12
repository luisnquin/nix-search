package nix_search_test

import (
	"context"
	"testing"

	"github.com/luisnquin/nix-search/internal/config"
	nix_search "github.com/luisnquin/nix-search/internal/nix/search"
)

func TestFlakeOptionsSmoke(t *testing.T) {
	ctx, config := context.Background(), config.Load()

	client, err := nix_search.NewClient(ctx, config)
	if err != nil {
		t.Fatal(err)
	}

	options, err := client.SearchFlakeOptions(ctx, "wayland", 50)
	if err != nil {
		t.Fatal(err)
	}

	for _, option := range options {
		t.Log(option)
	}
}

func TestFlakePackagesSmoke(t *testing.T) {
	ctx, config := context.Background(), config.Load()

	client, err := nix_search.NewClient(ctx, config)
	if err != nil {
		t.Fatal(err)
	}

	options, err := client.SearchFlakePackages(ctx, "wayland", 50)
	if err != nil {
		t.Fatal(err)
	}

	for _, option := range options {
		t.Log(option)
	}
}
