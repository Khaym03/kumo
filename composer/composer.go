package composer

import (
	"github.com/Khaym03/kumo/config"
	"github.com/Khaym03/kumo/controller"
	"github.com/Khaym03/kumo/core"
	sqlite "github.com/Khaym03/kumo/db/sqlite"
	db "github.com/Khaym03/kumo/db/sqlite/gen"
	"github.com/Khaym03/kumo/pkg/browser"
	"github.com/Khaym03/kumo/proxy"
	sche "github.com/Khaym03/kumo/scheduler"
	_ "github.com/go-rod/stealth"
	log "github.com/sirupsen/logrus"
)

// The AppComposer struct acts as the central hub for dependency injection.
// It holds the application's configuration and knows how to create
// and connect all other components.
type AppComposer struct {
	conf       config.AppConfig
	browserFac *browser.BrowserFactory
	browsers   []browser.BrowserCreator
	config.Config
	proxies []proxy.Proxy
}

// const limitOfBrowserInstances = 1

func NewAppComposer() *AppComposer {
	conn := sqlite.NewSQLiteConn()
	queries := db.New(conn)

	conf := config.LoadKumoConfig()

	// p, err := proxy.NewWebshareProxyProvider().Download()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	p := []proxy.Proxy{}

	pm := proxy.NewConcurrentProxyManager(p)

	browsers := browser.CreateCreatorsFromConfig(conf.Browsers, pm)

	return &AppComposer{
		conf: config.NewAppConfig(
			conn,
			queries,
			config.NewTaskStates(queries),
		),
		browserFac: browser.NewFactory(browsers...),
		browsers:   browsers,
		proxies:    p,
	}
}

// ComposeKumo builds and returns a Kumo instance with all its dependencies.
func (ac *AppComposer) ComposeKumo() (*core.Kumo, error) {
	// fill the browser pull
	b, err := ac.browserFac.Get()
	if err != nil {
		log.Fatal(err)
	}

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
