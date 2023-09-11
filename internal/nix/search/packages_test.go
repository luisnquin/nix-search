package nix_search_test

import (
	"context"
	"testing"

	"github.com/luisnquin/nix-search/internal/config"
	nix_search "github.com/luisnquin/nix-search/internal/nix/search"
)

func TestPackagesSmoke(t *testing.T) {
	const TEST_CHANNEL = "latest-42-nixos-unstable"

	ctx, config := context.Background(), config.Load()

	client, err := nix_search.NewClient(ctx, config)
	if err != nil {
		t.Fatal(err)
	}

	options, err := client.SearchPackages(ctx, TEST_CHANNEL, "go", 50)
	if err != nil {
		t.Fatal(err)
	}

	for _, option := range options {
		t.Log(option.Name)
	}
}
