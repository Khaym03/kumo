package controller

import (
	"context"
	"fmt"

	"github.com/Khaym03/kumo/config"
	db "github.com/Khaym03/kumo/db/sqlite/gen"
)

// FetchTasks obtain all the task to be proccess
type TaskFetcher interface {
	FetchTasks(ctx context.Context) ([]db.Task, error)
}

func NewTaskFetcher(conf config.AppConfig) TaskFetcher {
	return &taskFetcher{
		AppConfig: conf,
	}
}

type taskFetcher struct {
	config.AppConfig
}

func (f *taskFetcher) FetchTasks(ctx context.Context) ([]db.Task, error) {
	tasks, err := f.Queries.ListPendingOrFailedTasks(ctx, db.ListPendingOrFailedTasksParams{
		StatusID:   f.TaskStatus.Pending().ID,
		StatusID_2: f.TaskStatus.Failed().ID,
	})

	if err != nil {
		return nil, fmt.Errorf("error getting tasks from the database: %w", err)
	}

	return tasks, nil
}
