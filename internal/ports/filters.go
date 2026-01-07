package ports

import "github.com/Khaym03/kumo/internal/pkg/types"

// RequestFilter defines a contract for filtering requests.
type RequestFilter interface {
	Filter(req []*types.Request) []*types.Request
}
