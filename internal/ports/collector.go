package ports

import "context"

// Collector is a single, independent unit of work. It represents a specific process,
// like scraping a webpage for product details or extracting information from an API
type Collector interface {
	Collect(context.Context) error
}

//  Designed to manage and organize multiple Collector instances
type CollectorRegistry interface {
	RegisterCollector(Collector)
	Collectors() []Collector
}
