package collectors

import (
	"context"
	"net/url"

	sche "github.com/Khaym03/kumo/scheduler"
)

// Collect all the URLs require for a given task
type Collector interface {
	Collect(context.Context, sche.Scheduler) error
	URLs() []url.URL
}
