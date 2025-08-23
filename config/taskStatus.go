package config

import (
	"context"
	"log"

	db "github.com/Khaym03/kumo/db/sqlite/gen"
	"github.com/Khaym03/kumo/ports"
)

type taskStatus [4]db.TaskStatus

func NewTaskStates(queries *db.Queries) ports.TaskStatus {
	status, err := queries.GetTaskStatus(context.Background())
	if err != nil {
		log.Fatalf("Error at making [TaskStatus] %v", err)
	}

	if len(status) != 4 {
		log.Fatal("expect 4 status")
	}

	var taskStates taskStatus

	copy(taskStates[:], status)

	return &taskStates
}

func (ts *taskStatus) Pending() db.TaskStatus {
	return ts[0]
}

func (ts *taskStatus) InProgress() db.TaskStatus {
	return ts[1]
}

func (ts *taskStatus) Completed() db.TaskStatus {
	return ts[2]
}

func (ts *taskStatus) Failed() db.TaskStatus {
	return ts[3]
}
