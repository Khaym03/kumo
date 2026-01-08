package core

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/Khaym03/kumo/internal/pkg/types"
	"github.com/Khaym03/kumo/internal/ports"

	_ "github.com/PuerkitoBio/goquery"
	_ "golang.org/x/time/rate"
)

// KumoEngine is the core orchestration engine for the web crawler.
type KumoEngine struct {
	ctx        context.Context
	collectors map[string]ports.Collector
	resources *ResourcePool
	dispatcher ports.Dispatcher
	workersCount int
}

// NewKumoEngine creates and initializes a new Kumo engine.
func NewKumoEngine(
	ctx context.Context,
	rp *ResourcePool,
	dispather ports.Dispatcher,
	cs ...ports.Collector,
) *KumoEngine {
	m := map[string]ports.Collector{}
	for _, c := range cs {
		m[c.String()] = c
	}

	return &KumoEngine{
		ctx:          ctx,
		resources: rp,
		collectors:   m,
		dispatcher: dispather,
		workersCount: rp.Pages.Size(),
	}
}

// Run starts the crawling process, either resuming from a previous run or starting a new one.
func (k *KumoEngine) Run(initialReqs ...*types.Request) error {
	log.Println("Starting Kumo engine...")

	pending, err := k.dispatcher.LoadSavedState()
	if err != nil {
		return err
	}

	allRequests := append(pending, initialReqs...)
    if len(allRequests) == 0 {
        log.Warn("No requests to process. Exiting.")
        return nil
    }
    
    k.dispatcher.Dispatch(allRequests...)

	workersWaitGrp := sync.WaitGroup{}
	for i := range k.workersCount {
		workersWaitGrp.Go(func() {
			k.worker(i+1)
		})
	}

	k.dispatcher.Shutdown()
	workersWaitGrp.Wait()

	return nil
}


// Shutdown gracefully shuts down the engine and its dependencies.
func (k *KumoEngine) Shutdown() error {
	return k.resources.Browsers.Close()
}

func (k *KumoEngine) worker(id int) {
	for {
		// Pull blocks until a task is available or the queue is closed
		req, ok := k.dispatcher.Pull(k.ctx)
		if !ok {
			log.Debugf("Worker %d: no more tasks, exiting", id)
			return
		}
		
		k.executeTask(req)
	}
}

func (k *KumoEngine) executeTask(req *types.Request) {
    var processErr error
    
    defer func() {
        k.dispatcher.Finish(req, processErr)
    }()

    collector, ok := k.collectors[req.Collector]
    if !ok {
        log.Errorf("Collector %s not found for URL %s", req.Collector, req.URL)
        return
    }

    page, err := k.resources.Pages.Get()
    if err != nil {
        processErr = err
        return
    }
    defer k.resources.Pages.Put(page)

    // Execute the actual scraping logic
    processErr = collector.ProcessPage(k.ctx, page, req, k.dispatcher)
}

