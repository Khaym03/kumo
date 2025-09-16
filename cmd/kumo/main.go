package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Khaym03/kumo/internal/adapters/config"
	"github.com/Khaym03/kumo/internal/adapters/pagepool"
	"github.com/Khaym03/kumo/internal/adapters/storage"
	"github.com/Khaym03/kumo/internal/core"
	"github.com/Khaym03/kumo/internal/pkg/browser"
	"github.com/Khaym03/kumo/internal/pkg/proxy"
	"github.com/Khaym03/kumo/internal/pkg/types"
	"github.com/Khaym03/kumo/internal/ports"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	log.Info("loading config...")
	conf := config.LoadKumoConfig()

	log.Println(conf.String())
	pm := proxy.NewConcurrentProxyManager([]proxy.Proxy{})

	creators := browser.CreateCreatorsFromConfig(conf.Browsers, pm)
	browserPool := browser.NewPool(creators...)
	pp := pagepool.NewPagePool(browserPool, conf.NumOfPagesPerBrowser)

	dbConn, err := storage.NewBadgerDB(conf.StorageDir, conf.AllowBadgerLogger)
	if err != nil {
		log.Fatal(err)
	}

	db := storage.NewBadgerDBStore(dbConn)

	collectors := []ports.Collector{
		// concrete collectors
	}

	kumo := core.NewKumoEngine(
		ctx,
		browserPool,
		pp,
		db,
		db,
		collectors...,
	)
	defer func() {
		err = kumo.Shutdown()
		if err != nil {
			log.Warn(err)
		}
	}()

	go func() {
		// Must pass an initial ...*types.Request
		kumo.Run(&types.Request{})
		cancel()
	}()

	// The select statement waits for either an OS signal or a completion signal.
	select {
	case sig := <-sigChan:
		log.Infof("Interruption signal received (%v).", sig)
		cancel()
	case <-ctx.Done():
		log.Info("The scraping process completed its normal execution.")
	}

	log.Info("Application shutting down gracefully.")

}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	log.SetFormatter(&log.TextFormatter{ForceColors: true})
}
