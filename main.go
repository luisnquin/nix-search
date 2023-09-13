package main

import (
	"context"
	"log"

	"github.com/luisnquin/nix-search/internal/app"
	"github.com/luisnquin/nix-search/internal/config"
)

func main() {
	ctx := context.Background()
	appConfig := config.Load()

	app, err := app.New(appConfig)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
