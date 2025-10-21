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
	queue      chan *types.Request
	requestWg  *sync.WaitGroup
	mu         *sync.Mutex
	collectors map[string]ports.Collector
	ports.BrowserPool
	ports.PagePool
	fs           ports.FileStorage
	reqStore     ports.PersistenceStore
	reqFilters   []ports.RequestFilter
	workersCount int
}

// NewKumoEngine creates and initializes a new Kumo engine.
func NewKumoEngine(
	ctx context.Context,
	bp ports.BrowserPool,
	pp ports.PagePool,
	rs ports.PersistenceStore,
	fs ports.FileStorage,
	rf []ports.RequestFilter,
	cs ...ports.Collector,
) *KumoEngine {
	m := map[string]ports.Collector{}
	for _, c := range cs {
		m[c.String()] = c
	}

	return &KumoEngine{
		ctx:          ctx,
		queue:        make(chan *types.Request, 500),
		requestWg:    new(sync.WaitGroup),
		mu:           &sync.Mutex{},
		BrowserPool:  bp,
		PagePool:     pp,
		collectors:   m,
		reqStore:     rs,
		fs:           fs,
		reqFilters:   rf,
		workersCount: pp.Size(),
	}
}

// Run starts the crawling process, either resuming from a previous run or starting a new one.
func (k *KumoEngine) Run(initialReqs ...*types.Request) error {
	log.Println("Starting Kumo engine...")

	pending, err := k.reqStore.LoadPending()
	if err != nil {
		return err
	}

	if len(pending) == 0 && len(initialReqs) == 0 {
		log.Info("No initial requests and no pending requests. Exiting.")
		return nil
	}

	if len(pending) > 0 {
		log.Infof("Resuming from previous crawl with %d pending requests.", len(pending))
		k.Enqueue(pending...)
	} else {
		log.Info("Starting new crawl with initial requests.")
		k.Enqueue(initialReqs...)
	}

	workersWaitGrp := sync.WaitGroup{}
	for i := 1; i <= k.workersCount; i++ {
		workersWaitGrp.Go(func() {
			k.worker(k.ctx, i)
		})
	}

	k.requestWg.Wait()

	close(k.queue)

	workersWaitGrp.Wait()

	return nil
}

// Enqueue adds new requests to the processing queue after applying filters.
func (k *KumoEngine) Enqueue(r ...*types.Request) {
	k.mu.Lock()
	defer k.mu.Unlock()

	requestsToSave := k.filterRequests(r)

	if len(requestsToSave) == 0 {
		return
	}

	if err := k.reqStore.SavePending(requestsToSave...); err != nil {
		log.Printf("Error saving requests to store: %v", err)
		return
	}

	k.requestWg.Add(len(requestsToSave))
	go k.dispatchRequests(requestsToSave)
}

// Shutdown gracefully shuts down the engine and its dependencies.
func (k *KumoEngine) Shutdown() error {
	return k.BrowserPool.Close()
}

// filterRequests applies all registered filters to a slice of requests.
func (k *KumoEngine) filterRequests(reqs []*types.Request) []*types.Request {
	filteredReq := []*types.Request{}
	for _, req := range reqs {
		shouldFilter := false
		for _, filter := range k.reqFilters {
			filtered, err := filter.Filter(req)
			if err != nil {
				log.Warnf("Error applying filter: %v", err)
				continue
			}
			if filtered {
				shouldFilter = true
				break
			}
		}
		if !shouldFilter {
			filteredReq = append(filteredReq, req)
		}
	}
	return filteredReq
}

// worker fetches requests from the queue and processes them.
func (k *KumoEngine) worker(ctx context.Context, id int) {
	for {
		select {
		case req, ok := <-k.queue:
			if !ok {
				log.Printf("Worker %d queue closed. Exiting...", id)
				return
			}
			log.Printf("Worker %d processing URL: %s", id, req.URL)

			collector, ok := k.collectors[req.Collector]
			if !ok {
				log.Printf("No collector found for: %s", req.Collector)
				k.requestWg.Done()
				continue
			}

			page, err := k.PagePool.Get()
			if err != nil {
				log.Println("Error getting page from pool:", err)
				k.requestWg.Done()
				continue
			}

			err = collector.ProcessPage(ctx, page, req, k, k.fs)
			if err != nil {
				log.Warn("Error processing page:", err)
			} else {
				if err := k.reqStore.SaveCompleted(req); err != nil {
					log.Printf("Error saving completed request: %v", err)
				}
				if err := k.reqStore.RemoveFromPending(req); err != nil {
					log.Printf("Error removing from pending: %v", err)
				}
			}

			k.PagePool.Put(page)
			k.requestWg.Done()
		case <-ctx.Done():
			log.Printf("Worker %d shutting down gracefully.", id)
			return
		}
	}
}

// dispatchRequests sends requests to the queue in a separate goroutine.
func (k *KumoEngine) dispatchRequests(requests []*types.Request) {
	for _, req := range requests {
		k.queue <- req
	}
}
