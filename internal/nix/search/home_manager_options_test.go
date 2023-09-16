package nix_search_test

import (
	"context"
	"testing"

	"github.com/luisnquin/nix-search/internal/config"
	nix_search "github.com/luisnquin/nix-search/internal/nix/search"
)

func TestHomeManagerTest(t *testing.T) {
	appConfig, err := config.Load(true)
	if err != nil {
		t.Fatal(err)
	}

	client := nix_search.NewClient(appConfig)
	ctx := context.Background()

	options, err := client.SearchHomeManagerOptions(ctx, "programs.go")
	if err != nil {
		t.Fatal(err)
	}

	client.SearchHomeManagerOptions(ctx, "programs.bat")

	client.SearchHomeManagerOptions(ctx, "programs.wayland")
	client.SearchHomeManagerOptions(ctx, "programs.vscode")

	for _, option := range options {
		t.Log(option.Title)
	}
}
