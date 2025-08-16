package config

import (
	"context"

	db "github.com/Khaym03/kumo/db/sqlite/gen"
	"github.com/sirupsen/logrus"
)

type AppConfig struct {
	Queries    *db.Queries
	Logger     *logrus.Logger
	TaskStatus *taskStatus
}

func NewAppConfig(ctx context.Context, queries *db.Queries) AppConfig {
	return AppConfig{
		Queries:    queries,
		TaskStatus: NewTaskStates(queries),
		Logger:     logrus.New(),
	}
}
