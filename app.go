package main

import (
	"context"
	"fmt"
	"sync"

	adapters "github.com/Khaym03/kumo/internal/adapters/collector"
	"github.com/Khaym03/kumo/internal/adapters/config"
	"github.com/Khaym03/kumo/internal/adapters/filter"
	"github.com/Khaym03/kumo/internal/adapters/pagepool"
	"github.com/Khaym03/kumo/internal/adapters/storage"
	"github.com/Khaym03/kumo/internal/core"
	"github.com/Khaym03/kumo/internal/pkg/browser"
	"github.com/Khaym03/kumo/internal/pkg/proxy"
	"github.com/Khaym03/kumo/internal/pkg/types"
	"github.com/Khaym03/kumo/internal/ports"
	log "github.com/sirupsen/logrus"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {

	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

var (
	// runContext context.Context
	cancelFunc context.CancelFunc
	mu         sync.Mutex // Mutex to protect cancelFunc
)

func (a *App) RunKumo(conf KumoConfig) {
	fmt.Println(conf)

	// Lock to check the state and set the new context safely
	mu.Lock()
	if cancelFunc != nil {
		mu.Unlock()
		log.Println("A task is already running.")
		return
	}
	ctx, cancel := context.WithCancel(a.ctx)
	cancelFunc = cancel
	mu.Unlock() // Unlock after setting the state

	// --- INITIALIZE SERVICES ---
	// Initialize the proxy manager, browser pool, and page pool.
	// You can configure different proxy managers here.
	pm := proxy.NewConcurrentProxyManager([]proxy.Proxy{})
	creators := browser.CreateCreatorsFromConfig(conf.Browsers, pm)
	browserPool := browser.NewPool(creators...)
	pp := pagepool.NewPagePool(browserPool, 2)

	// Initialize the persistence store. The application will use this for
	// caching and tracking completed/pending requests.
	dbConn, err := storage.NewBadgerDB("./temp/storage", false)
	if err != nil {
		log.Fatal(err)
	}
	// Defer closing the database connection.
	defer dbConn.Close()
	db := storage.NewBadgerDBStore(dbConn)

	// --- DEFINE COLLECTORS AND INITIAL REQUESTS ---
	// Register your concrete collectors here. Collectors define the
	// scraping logic for different types of URLs.
	collectors := []ports.Collector{
		&adapters.AnugaPhase1Collector{},
		&adapters.AnugaPhase2Collector{},
	}

	// Define the starting URLs for the crawl.
	initialRequests := initialsURL(collectors[0].String())

	// Define the logic to skip types.Request
	requestFilters := []ports.RequestFilter{
		filter.NewIsCompletedFilter(db),
	}

	// --- START THE ENGINE ---
	kumo := core.NewKumoEngine(
		ctx,
		browserPool,
		pp,
		db,
		db,
		requestFilters,
		collectors...,
	)

	go func() {
		// --- RUN THE ENGINE AND WAIT FOR IT TO FINISH ---
		// This function will block until the engine stops.
		if err := kumo.Run(initialRequests...); err != nil {
			// Log the error but don't panic.
			if err == context.Canceled {
				log.Info("Kumo engine was canceled gracefully.")
			} else {
				log.Errorf("Kumo engine returned an error: %v", err)
			}
		}
		cancel()
	}()

	<-ctx.Done()

	// --- CLEANUP ---
	// Shutdown the engine and reset the state
	if err := kumo.Shutdown(); err != nil {
		log.Error(err)
	}
	// Lock again to safely reset the state
	mu.Lock()
	cancelFunc = nil
	mu.Unlock()
}

func (a *App) CancelKumo() {
	mu.Lock()
	defer mu.Unlock()

	if cancelFunc != nil {
		cancelFunc() // Signal the running task to stop.
		fmt.Println("Cancellation signal sent.")
		return
	}

	fmt.Println("No task is currently running.")

}

type KumoConfig struct {
	Root     string                 `json:"root"`
	Browsers []config.BrowserConfig `json:"browsers"`
}

func initialsURL(collectorName string) []*types.Request {
	var reqs []*types.Request
	for i := 0; i <= 7920; i += 20 {
		reqs = append(reqs, &types.Request{
			URL:       fmt.Sprintf(`https://www.anuga.com/anuga-exhibitors/list-of-exhibitors/?route=aussteller/blaettern&fw_ajax=5&paginatevalues={"stichwort":"","suchart":"alle"}&start=%d&dat=351214`, i),
			Collector: collectorName,
		})
	}

	return reqs
}
