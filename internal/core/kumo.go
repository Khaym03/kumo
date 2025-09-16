package core

import (
	"context"
	"sync"

	"github.com/Khaym03/kumo/internal/pkg/types"
	"github.com/Khaym03/kumo/internal/ports"
	log "github.com/sirupsen/logrus"

	_ "github.com/PuerkitoBio/goquery"
	_ "golang.org/x/time/rate"
)

type KumoEngine struct {
	ctx        context.Context
	queue      chan *types.Request
	wg         *sync.WaitGroup
	mu         *sync.Mutex
	collectors map[string]ports.Collector
	ports.BrowserPool
	ports.PagePool
	fs           ports.FileStorage
	reqStore     ports.PersistenceStore
	workersCount int
}

func NewKumoEngine(
	ctx context.Context,
	bp ports.BrowserPool,
	pp ports.PagePool,
	rs ports.PersistenceStore,
	fs ports.FileStorage,
	cs ...ports.Collector,
) *KumoEngine {
	m := map[string]ports.Collector{}
	for _, c := range cs {
		m[c.String()] = c
	}

	return &KumoEngine{
		ctx:          ctx,
		queue:        make(chan *types.Request, 500),
		wg:           new(sync.WaitGroup),
		mu:           &sync.Mutex{},
		BrowserPool:  bp,
		PagePool:     pp,
		collectors:   m,
		reqStore:     rs,
		fs:           fs,
		workersCount: pp.Size(),
	}
}

func (k *KumoEngine) Run(initialReqs ...*types.Request) error {
	log.Println("Starting Kumo engine...")

	// Phase 1: Check for pending requests and enqueue them.
	pending, err := k.reqStore.LoadPending()
	if err != nil {
		return err
	}

	if len(pending) > 0 {
		log.Infof("Resuming from previous crawl with %d pending requests.", len(pending))
		k.Enqueue(pending...)
	} else {
		// Only enqueue initial requests if there are no pending requests from a previous run.
		if len(initialReqs) > 0 {
			log.Info("Starting new crawl with initial requests.")
			k.Enqueue(initialReqs...)
		} else {
			log.Info("No initial requests and no pending requests. Exiting.")
			return nil
		}
	}

	// Start workers first
	workersWaitGrp := sync.WaitGroup{}
	workersWaitGrp.Add(k.workersCount)
	for i := 1; i <= k.workersCount; i++ {
		go k.worker(k.ctx, &workersWaitGrp, i)
	}

	// Wait for all tasks to be completed.
	k.wg.Wait()

	// Signal workers to stop by closing the queue.
	close(k.queue)

	// Wait for all workers to shut down cleanly.
	workersWaitGrp.Wait()

	return nil
}

func (k *KumoEngine) Enqueue(r ...*types.Request) {
	k.mu.Lock()
	requestsToSave := []*types.Request{}
	for _, req := range r {
		isCompleted, err := k.reqStore.IsCompleted(req.URL)
		if err != nil {
			log.Printf("Error checking if request is completed: %v", err)
			continue
		}
		if !isCompleted {
			requestsToSave = append(requestsToSave, req)
		}
	}
	k.mu.Unlock()

	if len(requestsToSave) == 0 {
		return
	}

	if err := k.reqStore.SavePending(requestsToSave...); err != nil {
		log.Printf("Error saving requests to store: %v", err)
		return
	}

	k.wg.Add(len(requestsToSave))
	go k.dispatchRequests(requestsToSave)
}

// dispatchRequests sends requests to the queue in a separate goroutine.
// This is used for bulk operations (like initial requests) to avoid blocking the main thread.
func (k *KumoEngine) dispatchRequests(requests []*types.Request) {
	for _, req := range requests {
		k.queue <- req
	}
}

func (k *KumoEngine) worker(ctx context.Context, wg *sync.WaitGroup, id int) {
	defer wg.Done()

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
				k.wg.Done()
				continue
			}

			page, err := k.PagePool.Get()
			if err != nil {
				log.Println("Error getting page from pool:", err)
				k.wg.Done()
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
			k.wg.Done()
		case <-ctx.Done():
			log.Printf("Worker %d shutting down gracefully.", id)
			return
		}
	}
}

func (k *KumoEngine) Shutdown() error {
	err := k.reqStore.Close()
	if err != nil {
		log.Error(err)
	}
	return k.BrowserPool.Close()
}
