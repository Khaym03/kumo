package filter

import (
	"github.com/Khaym03/kumo/internal/pkg/types"
	"github.com/Khaym03/kumo/internal/ports"
)

type filterComposite struct {
	filters []ports.RequestFilter
}

func NewFilterComposite(filters ...ports.RequestFilter) ports.RequestFilter {
	return &filterComposite{
		filters: filters,
	}
}

// Filter executes all registered filters sequentially.
func (fc *filterComposite) Filter(reqs []*types.Request) []*types.Request {
	filteredReqs := reqs

	for _, f := range fc.filters {
		filteredReqs = f.Filter(filteredReqs)
	}

	return filteredReqs
}
