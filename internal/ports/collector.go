package ports

import (
	"context"

	"github.com/Khaym03/kumo/internal/pkg/types"
	"github.com/go-rod/rod"
)

// Collector is a single, independent unit of work. It represents a specific process,
// like scraping a webpage for product details or extracting information from an API
// ports/collector.go
type Collector interface {
	ProcessPage(
		ctx context.Context,
		page *rod.Page,
		req *types.Request,
		queue Enqueuer,
		fs FileStorage,
	) error

	String() string
}

// Designed to manage and organize multiple Collector instances
type CollectorRegistry interface {
	RegisterCollector(Collector)
	Collectors() []Collector
}

type Enqueuer interface {
	Enqueue(reqs ...*types.Request)
}
