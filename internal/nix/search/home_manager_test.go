package nix_search_test

import (
	"context"
	"testing"

	"github.com/luisnquin/nix-search/internal/config"
	nix_search "github.com/luisnquin/nix-search/internal/nix/search"
)

func TestHomeManagerTest(t *testing.T) {
	ctx, config := context.Background(), config.Load()

	client, err := nix_search.NewClient(ctx, config)
	if err != nil {
		t.Fatal(err)
	}

	options := client.SearchHomeManagerOptions("programs.go")

	for _, option := range options {
		t.Log(option.Title)
	}
}
