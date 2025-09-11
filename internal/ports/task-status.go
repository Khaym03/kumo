package ports

import db "github.com/Khaym03/kumo/db/sqlite/gen"

type TaskStatus interface {
	Pending() db.TaskStatus
	InProgress() db.TaskStatus
	Completed() db.TaskStatus
	Failed() db.TaskStatus
}
