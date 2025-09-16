package filter

import (
	"github.com/Khaym03/kumo/internal/pkg/types"
	"github.com/Khaym03/kumo/internal/ports"
)

// IsCompletedFilter checks if a request has already been completed.
type IsCompletedFilter struct {
	checker ports.CompletionChecker
}

func NewIsCompletedFilter(c ports.CompletionChecker) ports.RequestFilter {
	return &IsCompletedFilter{checker: c}
}

func (f *IsCompletedFilter) Filter(req *types.Request) (bool, error) {
	isCompleted, err := f.checker.IsCompleted(req.URL)
	if err != nil {
		return false, err
	}
	return isCompleted, nil // returns true if it should be filtered (skipped)
}
