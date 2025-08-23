package config

import (
	"database/sql"

	db "github.com/Khaym03/kumo/db/sqlite/gen"
	"github.com/Khaym03/kumo/ports"
	"github.com/sirupsen/logrus"
)

type AppConfig struct {
	Queries    *db.Queries
	Logger     *logrus.Logger
	TaskStatus ports.TaskStatus
	DB         *sql.DB
}

func NewAppConfig(
	DB *sql.DB,
	queries *db.Queries,
	logger *logrus.Logger,
	ts ports.TaskStatus,
) AppConfig {
	return AppConfig{
		Queries:    queries,
		TaskStatus: ts,
		Logger:     logger,
		DB:         DB,
	}
}
