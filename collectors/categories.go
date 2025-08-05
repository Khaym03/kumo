package collectors

import (
	"context"
	"log"
	"net/url"
	"time"

	sche "github.com/Khaym03/kumo/scheduler"
)

type CategoriesCollector struct {
	urls []url.URL
}

func (cc *CategoriesCollector) Collect(ctx context.Context, scheduler sche.Scheduler) error {
	p, err := scheduler.Get()
	if err != nil {
		return err
	}
	defer scheduler.Put(p)

	// logic
	log.Println("doing some job...")
	time.Sleep(2 * time.Second)

	return nil
}

func (cc *CategoriesCollector) URLs() []url.URL {
	return cc.urls
}
