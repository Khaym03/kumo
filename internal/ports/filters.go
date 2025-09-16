package ports

import "github.com/Khaym03/kumo/internal/pkg/types"

// RequestFilter defines a contract for filtering requests.
// It returns a boolean to indicate if the request should be filtered (true) or not (false),
// and an error if something goes wrong.
type RequestFilter interface {
	Filter(req *types.Request) (bool, error)
}
