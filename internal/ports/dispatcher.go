package ports

import (
	"context"

	"github.com/Khaym03/kumo/internal/pkg/types"
)

type Dispatcher interface {
	// return pending reqs if any
	LoadSavedState() ([]*types.Request, error)
	Dispatch(reqs ...*types.Request)
	Finish(req *types.Request, err error)
	Pull(ctx context.Context) (*types.Request, bool)
	Shutdown()
}