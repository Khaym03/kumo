package core

import (
	"context"
	"log"

	"github.com/Khaym03/kumo/collectors"
	"github.com/Khaym03/kumo/config"
	"github.com/Khaym03/kumo/controller"
	sche "github.com/Khaym03/kumo/scheduler"
	"github.com/go-rod/rod"

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
}

func NewKumo(
	browser *rod.Browser,
	scheduler sche.Scheduler,
	registry *collectors.CollectorRegistry,
	reconciler controller.Reconciler,
	appConfig config.AppConfig,
) *Kumo {
	return &Kumo{
		browser:           browser,
		scheduler:         scheduler,
		CollectorRegistry: registry,
		Reconciler:        reconciler,
		AppConfig:         appConfig,
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
}
