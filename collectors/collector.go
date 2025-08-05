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
