package main

import (
	"context"
	"log"

	"github.com/luisnquin/flaggy"
	"github.com/luisnquin/nix-search/internal"
	"github.com/luisnquin/nix-search/internal/config"
	"github.com/luisnquin/nix-search/internal/gui"
)

var version string = internal.DEFAULT_VERSION

func main() {
	flaggy.SetName(internal.PROGRAM_NAME)
	flaggy.SetDescription(internal.PROGRAM_DESCRIPTION)
	flaggy.SetVersion(version)
	flaggy.DefaultParser.SetHelpTemplate(internal.GetHelpTemplate())
	flaggy.Parse()

	ctx := context.Background()
	appConfig := config.Load()

	switch {
	default:
		gui, err := gui.New(appConfig)
		if err != nil {
			log.Fatal(err)
		}

		if err := gui.Run(ctx); err != nil {
			log.Fatal(err)
		}
	}
}
