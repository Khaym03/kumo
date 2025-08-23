package ports

import (
	"context"

	db "github.com/Khaym03/kumo/db/sqlite/gen"
)

// Exclude tasks that have failed more
type PendingOrFailedTaskLister interface {
	ListPendingOrFailedTasks(ctx context.Context) ([]db.Task, error)
}
