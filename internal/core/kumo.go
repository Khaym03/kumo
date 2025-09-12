package core

import (
	"context"

	"log"
	"sync"

	"github.com/Khaym03/kumo/internal/pkg/types"
	"github.com/Khaym03/kumo/internal/ports"

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
	dataSink ports.DataSink
}

func NewKumoEngine(
	bp ports.BrowserPool,
	pp ports.PagePool,
	ds ports.DataSink,
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
		dataSink:    ds,
		collectors:  m,
	}
}

func (k *KumoEngine) Run() error {
	log.Println("Starting Kumo engine...")

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

			err = collector.ProcessPage(ctx, page, req, k, k.dataSink)
			if err != nil {
				log.Println("Error processing page:", err)
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
	k.wg.Add(len(r))
	k.queue.Enqueue(r...)
}

func (k *KumoEngine) Shutdown() error {
	return k.BrowserPool.Close()
}
