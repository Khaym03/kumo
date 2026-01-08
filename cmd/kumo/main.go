package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Khaym03/kumo/internal/adapters/config"
	"github.com/Khaym03/kumo/internal/adapters/filter"
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

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
}

// main orchestrates the scraping process, handling setup, execution, and graceful shutdown.
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// --- CONFIGURATION ---
	// Load the crawler's configuration from environment variables or a file.
	conf := config.LoadKumoConfig()
	log.Info("Loading config...")
	log.Println(conf.String())

	// --- INITIALIZE SERVICES ---
	// Initialize the proxy manager, browser pool, and page pool.
	// You can configure different proxy managers here.
	pm := proxy.NewConcurrentProxyManager([]proxy.Proxy{})
	creators := browser.CreateCreatorsFromConfig(conf.Browsers, pm)
	browserPool := browser.NewPool(creators...)
	pp := pagepool.NewPagePool(browserPool, conf.NumOfPagesPerBrowser)

	// Initialize the persistence store. The application will use this for
	// caching and tracking completed/pending requests.
	dbConn, err := storage.NewBadgerDB(conf.StorageDir, conf.AllowBadgerLogger)
	if err != nil {
		log.Fatal(err)
	}
	// Defer closing the database connection.
	defer dbConn.Close()

	db := storage.NewBadgerDBStore(dbConn)
	defer db.Close()

	// --- DEFINE COLLECTORS AND INITIAL REQUESTS ---
	// Register your concrete collectors here. Collectors define the
	// scraping logic for different types of URLs.
	collectors := []ports.Collector{
		// &adapters.AnugaPhase1Collector{},
		// &adapters.AnugaPhase2Collector{},
	}

	// Define the starting URLs for the crawl.
	initialRequests := []*types.Request{
		// &types.Request{URL: "https://example.com", Collector: "collector-name"},
	}

	// Define the logic to skip types.Request
	requestFilters := filter.NewFilterComposite(
		filter.NewIsCompletedFilter(db),
	)

	// --- START THE ENGINE ---
	kumo := core.NewKumoEngine(
		ctx,
		core.NewResourcePool(browserPool, pp),
		core.NewDispatcher(requestFilters, db),
		collectors...,
	)

	// Use a separate goroutine to run the engine so main can listen for signals.
	go func() {
		if err := kumo.Run(initialRequests...); err != nil {
			log.Errorf("Kumo engine returned an error: %v", err)
		}
		cancel() // Signal completion or error to the main goroutine
	}()

	// --- GRACEFUL SHUTDOWN HANDLING ---
	// Listen for OS signals to gracefully shut down the application.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		log.Infof("Interruption signal received (%v). Shutting down...", sig)
		// Cancel the context to signal all goroutines to stop.
	case <-ctx.Done():
		log.Info("The scraping process completed its normal execution.")
	}

	// The deferred cancel() and dbConn.Close() will clean up resources.
	log.Info("Application shutting down gracefully.")
}
