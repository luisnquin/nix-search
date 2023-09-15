package main

import (
	"context"
	"log"

	"github.com/luisnquin/nix-search/internal/config"
	"github.com/luisnquin/nix-search/internal/gui"
)

func main() {
	ctx := context.Background()
	appConfig := config.Load()

	gui, err := gui.New(appConfig)
	if err !=: nil {
		log.Fatal(err)
	}

	if err := gui.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
