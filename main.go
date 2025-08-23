package main

import (
	"log"

	"github.com/Khaym03/kumo/composer"
)

func main() {
	appComposer := composer.NewAppComposer()
	kumo, err := appComposer.ComposeKumo()
	if err != nil {
		log.Fatalf("Error al componer la aplicación: %v", err)
	}
	defer kumo.Shutdown()

	kumo.Run()
}
