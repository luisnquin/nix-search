package nix_search_test

import (
	"context"
	"testing"

	"github.com/luisnquin/nix-search/internal/config"
	nix_search "github.com/luisnquin/nix-search/internal/nix/search"
)

func TestNixOSOptionsSmoke(t *testing.T) {
	const TEST_CHANNEL = "latest-42-nixos-unstable"

	ctx, config := context.Background(), config.Load()

	client, err := nix_search.NewClient(ctx, config)
	if err != nil {
		t.Fatal(err)
	}

	options, err := client.SearchNixOSOptions(ctx, TEST_CHANNEL, "services.postgresql", 50)
	if err != nil {
		t.Fatal(err)
	}

	for _, option := range options {
		t.Log(option.Name)
	}
}
