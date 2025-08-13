package collectors

import (
	"context"
	"net/url"
)

// Collect all the URLs require for a given task
type Collector interface {
	Collect(context.Context) error
	URLs() []url.URL
}

type CollectorRegistry struct {
	Collectors []Collector
}

func NewCollectorRegistry() *CollectorRegistry {
	return &CollectorRegistry{
		Collectors: []Collector{},
	}
}

func (k *CollectorRegistry) RegisterCollector(c Collector) {
	k.Collectors = append(k.Collectors, c)
}
