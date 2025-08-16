package core

import (
	"context"
	"database/sql"
	"log"

	"github.com/Khaym03/kumo/collectors"
	"github.com/Khaym03/kumo/config"
	"github.com/Khaym03/kumo/controller"
	db "github.com/Khaym03/kumo/db/sqlite/gen"
	sche "github.com/Khaym03/kumo/scheduler"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"

	sqlite "github.com/Khaym03/kumo/db/sqlite"
	_ "github.com/PuerkitoBio/goquery"
	_ "golang.org/x/time/rate"
)

type Kumo struct {
	ctx       context.Context
	browser   *rod.Browser
	scheduler sche.Scheduler
	*collectors.CollectorRegistry
	controller.Reconciler
	config.AppConfig
	storage *sql.DB
}

func NewKumo() *Kumo {
	u := launcher.New().Leakless(false).MustLaunch()
	b := rod.New().ControlURL(u).Trace(true).MustConnect()

	conn := sqlite.NewSQLiteConn()

	conf := config.NewAppConfig(context.Background(), db.New(conn))

	return &Kumo{
		ctx:               context.Background(),
		browser:           b,
		scheduler:         sche.NewScheduler(b, 2),
		CollectorRegistry: collectors.NewCollectorRegistry(),
		Reconciler:        controller.NewStateReconciler(conf),
		AppConfig:         conf,
		storage:           conn,
	}
}

func (k *Kumo) Run() {
	err := k.Reconcile(k.ctx)
	if err != nil {
		k.Logger.Fatal(err)
	}

	for _, c := range k.Collectors {
		err := c.Collect(k.ctx)
		if err != nil {
			log.Println(err)
		}
	}
}

func (k Kumo) Shutdown() {
	k.browser.MustClose()
	// err := k.storage.Close()
	// if err != nil {
	// 	k.Logger.Error(err)
	// }
}
