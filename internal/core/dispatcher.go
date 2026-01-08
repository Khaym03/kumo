// core/req_manager.go
package core

import (
	"context"
	"sync"

	"github.com/Khaym03/kumo/internal/pkg/types"
	"github.com/Khaym03/kumo/internal/ports"
	log "github.com/sirupsen/logrus"
)

type Dispatcher struct {
	queue       chan *types.Request
	wg          sync.WaitGroup
	filter      ports.RequestFilter
	persistence ports.PersistenceStore
	closeOnce   sync.Once
}

func NewDispatcher(f ports.RequestFilter, p ports.PersistenceStore) *Dispatcher {
	return &Dispatcher{
		queue:       make(chan *types.Request, 500),
		filter:      f,
		persistence: p,
	}
}

func (d *Dispatcher) LoadSavedState() ([]*types.Request, error) {
    pending, err := d.persistence.LoadPending()
    if err != nil {
        return nil, err
    }
    
    // We don't dispatch them yet, we just return them to the Engine 
    // so it can decide how to handle them alongside new seeds.
    return pending, nil
}

func (rm *Dispatcher) Dispatch(reqs ...*types.Request) {
	filtered := rm.filter.Filter(reqs)
	if len(filtered) == 0 {
		return
	}

	rm.persistence.SavePending(filtered...)

	for _, r := range filtered {
		rm.wg.Add(1)
		rm.queue <- r
	}
}

func (rm *Dispatcher) Pull(ctx context.Context) (*types.Request, bool) {
	select {
	case <-ctx.Done():
		return nil, false
	case req, ok := <-rm.queue:
		return req, ok
	}
}

func (rm *Dispatcher) Finish(req *types.Request, err error) {
	defer rm.wg.Done()

	if err != nil {
		log.Errorf("Request %s failed, not marking as completed: %v", req.URL, err)
		return
	}

	// Only if err is nil, we mark it as finished in the DB
	if err := rm.persistence.SaveCompleted(req); err != nil {
		log.Errorf("failed to persist completed request: %v", err)
	}
	_ = rm.persistence.RemoveFromPending(req)
}

func (rm *Dispatcher) Shutdown() {
	rm.wg.Wait()
	rm.closeOnce.Do(func() {
		close(rm.queue)
	})
}