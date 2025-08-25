package main

import (
	"log"

	"github.com/Khaym03/kumo/composer"
	"github.com/joho/godotenv"
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

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
