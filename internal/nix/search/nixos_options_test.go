package nix_search_test

import (
	"context"
	"testing"

	"github.com/luisnquin/nix-search/internal/config"
	"github.com/luisnquin/nix-search/internal/nix"
	nix_search "github.com/luisnquin/nix-search/internal/nix/search"
)

func TestNixOSOptionsSmoke(t *testing.T) {
	appConfig, err := config.Load(true)
	if err != nil {
		t.Fatal(err)
	}

	client := nix_search.NewClient(appConfig)
	ctx := context.Background()

	channelStatus := nix.CHANNEL_STATUS_STABLE

	channel, found := appConfig.Internal.Nix.FindChannelWithStatus(channelStatus)
	if !found {
		t.Errorf("unable to find channel with %s status", channelStatus)
		t.FailNow()
	}

	options, err := client.SearchNixOSOptions(ctx, channel.Branch, "services.postgresql", 50)
	if err != nil {
		t.Fatal(err)
	}

	for _, option := range options {
		t.Log(option.Name)
	}
}
