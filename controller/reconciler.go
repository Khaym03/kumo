package controller

import (
	"context"

	"github.com/Khaym03/kumo/config"
	db "github.com/Khaym03/kumo/db/sqlite/gen"
)

// Responsable for continue the progress since the last use
type Reconciler interface {
	Reconcile(ctx context.Context) error
}

type StateReconciler struct {
	conf config.AppConfig
}

func NewStateReconciler(conf config.AppConfig) *StateReconciler {
	return &StateReconciler{
		conf: conf,
	}
}

func (r *StateReconciler) Reconcile(ctx context.Context) error {
	tasks, err := r.conf.Queries.ListTasksByStatusID(ctx, r.conf.TaskStatus.InProgress().ID)

	if err != nil {
		return err
	}

	if len(tasks) > 0 {
		r.conf.Logger.Infof("Found %d tasks in an 'IN_PROGRESS' state. Requeuing them...", len(tasks))
		for _, task := range tasks {
			err := r.conf.Queries.UpdateTaskStatus(ctx, db.UpdateTaskStatusParams{
				StatusID: r.conf.TaskStatus.Pending().ID,
				ID:       task.ID},
			)

			if err != nil {
				r.conf.Logger.Errorf("Failed to re-queue task %s: %v", task.ID, err)
			} else {
				r.conf.Logger.Infof("Task %s successfully re-queued.", task.ID)
			}
		}
	}
	return nil
}
