package filter

import (
	"github.com/Khaym03/kumo/internal/pkg/types"
	"github.com/Khaym03/kumo/internal/ports"
	log "github.com/sirupsen/logrus"
)

// IsCompletedFilter checks if a request has already been completed.
type IsCompletedFilter struct {
	checker ports.CompletionChecker
}

func NewIsCompletedFilter(c ports.CompletionChecker) ports.RequestFilter {
	return &IsCompletedFilter{checker: c}
}

func (f *IsCompletedFilter) Filter(reqs []*types.Request) []*types.Request {
	filteredReqs := make([]*types.Request, 0, len(reqs))

	for _, req := range reqs {
		isCompleted, err := f.checker.IsCompleted(req.URL)
		if err != nil {
			log.Warnf("fail to filter: %s", req.URL)
			continue
		}

		if isCompleted {
			continue
		}
		filteredReqs = append(filteredReqs, req)
	}

	return filteredReqs
}
