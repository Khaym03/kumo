package ports

import (
	"context"

	db "github.com/Khaym03/kumo/db/sqlite/gen"
)

// Interfaces for every unique behavior of [db.Queries]

type (
	// PendingOrFailedTaskLister lists tasks that are either pending or have failed,
	// excluding those that have failed more than 3 times.
	PendingOrFailedTaskLister interface {
		ListPendingOrFailedTasks(ctx context.Context) ([]db.Task, error)
	}

	// TaskAdder provides functionality to add a new task.
	TaskAdder interface {
		AddTask(ctx context.Context, arg db.AddTaskParams) error
	}

	// TaskFailer marks a task as failed and increments its retry count.
	TaskFailer interface {
		FailTask(ctx context.Context, id int64) error
	}

	// StatusIDGetter retrieves the ID of a status given its name.
	StatusIDGetter interface {
		GetStatusIDByName(ctx context.Context, name string) (int64, error)
	}

	// TaskStatusGetter obtains all task statuses.
	TaskStatusGetter interface {
		GetTaskStatus(ctx context.Context) ([]db.TaskStatus, error)
	}

	// TasksByStatusIDLister lists tasks filtered by a given status ID.
	TasksByStatusIDLister interface {
		ListTasksByStatusID(ctx context.Context, statusID int64) ([]db.Task, error)
	}

	// TaskStatusUpdater modifies the status of a task.
	TaskStatusUpdater interface {
		UpdateTaskStatus(ctx context.Context, arg db.UpdateTaskStatusParams) error
	}
)
