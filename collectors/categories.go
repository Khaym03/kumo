package collectors

import (
	"context"
	"log"
	"net/url"
	"time"

	sche "github.com/Khaym03/kumo/scheduler"
)

type CategoriesCollector struct {
	urls      []url.URL
	scheduler sche.Scheduler
}

func (cc *CategoriesCollector) Collect(ctx context.Context) error {
	p, err := cc.scheduler.Get()
	if err != nil {
		return err
	}
	defer cc.scheduler.Put(p)

	// logic
	log.Println("doing some job...")
	time.Sleep(2 * time.Second)

	return nil
}

func (cc *CategoriesCollector) URLs() []url.URL {
	return cc.urls
}

func NewCategoriesCollector(s sche.Scheduler) *CategoriesCollector {
	return &CategoriesCollector{scheduler: s}
}
