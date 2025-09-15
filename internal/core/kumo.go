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
	queue      types.RequestQueue
	wg         *sync.WaitGroup
	collectors map[string]ports.Collector
	ports.BrowserPool
	ports.PagePool
	reqStore ports.PersistenceStore
}

func NewKumoEngine(
	bp ports.BrowserPool,
	pp ports.PagePool,
	rs ports.PersistenceStore,
	cs ...ports.Collector,
) *KumoEngine {
	m := map[string]ports.Collector{}

	for _, c := range cs {
		m[c.Name()] = c
	}

	return &KumoEngine{
		ctx:         context.Background(),
		queue:       make(types.RequestQueue, 100),
		wg:          new(sync.WaitGroup),
		BrowserPool: bp,
		PagePool:    pp,
		collectors:  m,
		reqStore:    rs,
	}
}

func (k *KumoEngine) Run() error {
	log.Println("Starting Kumo engine...")

	pending, err := k.reqStore.LoadPending()
	if err != nil {
		return err
	}

	if len(pending) == 0 {
		log.Info("all URLs were processed")
		return nil
	}

	log.Infof("Pending requests %d", len(pending))

	k.Enqueue(pending...)

	go k.worker(k.ctx, 1)

	k.wg.Wait()    // Wait for all tasks to be completed.
	close(k.queue) // Signal the worker to stop by closing the queue.

	return k.Shutdown()
}

func (k *KumoEngine) worker(ctx context.Context, id int) {
	for {
		select {
		case req, ok := <-k.queue:
			if !ok {
				log.Println("Queue closed so return")
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

			err = collector.ProcessPage(ctx, page, req, k)
			if err != nil {
				log.Warn("Error processing page:", err)
			} else {
				// On successful completion, update the store.
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
			log.Printf("Worker %d shutting down.", id)
			return
		}
	}
}

func (k *KumoEngine) InitialReqs(r ...*types.Request) {
	k.Enqueue(r...)
}

func (k *KumoEngine) Enqueue(r ...*types.Request) {
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

	if len(requestsToSave) == 0 {
		return
	}

	if err := k.reqStore.SavePending(requestsToSave...); err != nil {
		log.Printf("Error saving requests to store: %v", err)
		return
	}

	k.wg.Add(len(requestsToSave))
	k.queue.Enqueue(requestsToSave...)
}

func (k *KumoEngine) Shutdown() error {
	err := k.reqStore.Close()
	if err != nil {
		log.Error(err)
	}
	return k.BrowserPool.Close()
}
