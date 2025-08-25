package composer

import (
	"github.com/Khaym03/kumo/config"
	"github.com/Khaym03/kumo/controller"
	"github.com/Khaym03/kumo/core"
	sqlite "github.com/Khaym03/kumo/db/sqlite"
	db "github.com/Khaym03/kumo/db/sqlite/gen"
	sche "github.com/Khaym03/kumo/scheduler"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	_ "github.com/go-rod/stealth"
	"github.com/sirupsen/logrus"
)

// The AppComposer struct acts as the central hub for dependency injection.
// It holds the application's configuration and knows how to create
// and connect all other components.
type AppComposer struct {
	conf config.AppConfig
}

func NewAppComposer() *AppComposer {
	conn := sqlite.NewSQLiteConn()
	queries := db.New(conn)

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	return &AppComposer{
		conf: config.NewAppConfig(
			conn,
			queries,
			logger,
			config.NewTaskStates(queries),
		),
	}
}

// ComposeKumo builds and returns a Kumo instance with all its dependencies.
func (ac *AppComposer) ComposeKumo() (*core.Kumo, error) {
	u := launcher.New().Leakless(false).MustLaunch()
	b := rod.New().ControlURL(u).Trace(true).MustConnect()

	scheduler := sche.NewScheduler(b, 2)
	reconciler := controller.NewStateReconciler(ac.conf)
	registry := controller.NewCollectorRegistry()

	kumo := core.NewKumo(
		b,
		scheduler,
		registry,
		reconciler,
		ac.conf,
	)

	return kumo, nil
}
