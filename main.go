package main

import (
	"context"

	"github.com/luisnquin/flaggy"
	"github.com/luisnquin/nix-search/internal"
	"github.com/luisnquin/nix-search/internal/config"
	"github.com/luisnquin/nix-search/internal/gui"
	"github.com/luisnquin/nix-search/internal/log"
	nix_search "github.com/luisnquin/nix-search/internal/nix/search"
)

var version string = internal.DEFAULT_VERSION

func main() {
	flaggy.SetName(internal.PROGRAM_NAME)
	flaggy.SetDescription(internal.PROGRAM_DESCRIPTION)
	flaggy.SetVersion(version)
	flaggy.DefaultParser.SetHelpTemplate(internal.GetHelpTemplate())
	flaggy.Parse()

	ctx := context.Background()

	appConfig, err := config.Load(false)
	if err != nil {
		log.Pretty.Fatal(err.Error())
	}

	logger, err := log.New(appConfig.LogFile)
	must(err)

	defer func() {
		must(logger.Close())
	}()

	nixClient := nix_search.NewClient(appConfig)

	switch {
	default:
		gui, err := gui.New(logger, appConfig, nixClient)
		if err != nil {
			logger.Err(err).Msg("unable to initialize program GUI")
			logger.Close()
			log.Pretty.Fatal(err.Error())
		}

		if err := gui.Run(ctx); err != nil {
			logger.Err(err).Msg("a problem happened during GUI execution")
			logger.Close()
			log.Pretty.Fatal(err.Error())
		}
	}
}

func must(err error) {
	if err != nil {
		log.Pretty.Fatal(err.Error())
	}
}
