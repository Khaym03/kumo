package config

import (
	"context"
	"log"

	db "github.com/Khaym03/kumo/db/sqlite/gen"
)

type TaskStatus [4]db.TaskStatus

func NewTaskStates(queries *db.Queries) *TaskStatus {
	status, err := queries.GetTaskStatus(context.Background())
	if err != nil {
		log.Fatalf("Error at making [TaskStatus] %v", err)
	}

	if len(status) != 4 {
		log.Fatal("expect 4 status")
	}

	var taskStates TaskStatus

	copy(taskStates[:], status)

	return &taskStates
}

func (ts *TaskStatus) Pending() db.TaskStatus {
	return ts[0]
}

func (ts *TaskStatus) InProgress() db.TaskStatus {
	return ts[1]
}

func (ts *TaskStatus) Completed() db.TaskStatus {
	return ts[2]
}

func (ts *TaskStatus) Failed() db.TaskStatus {
	return ts[3]
}
