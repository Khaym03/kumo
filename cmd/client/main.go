package main

import (
	"log"
	"time"

	"github.com/Khaym03/kumo/composer"
	"github.com/joho/godotenv"
)

func main() {
	ac := composer.NewAppComposer()

	browser, err := ac.ComposeRemoteBrowser()
	if err != nil {
		log.Fatal(err)
	}

	defer browser.MustClose()

	// Open a new page and interact with it
	page := browser.MustPage("https://google.com")
	page.MustWaitLoad()
	title := page.MustInfo().Title
	log.Printf("Page title: %s\n", title)

	time.Sleep(time.Second * 10)
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}
