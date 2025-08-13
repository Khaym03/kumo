package core

import (
	"context"
	"log"

	"github.com/Khaym03/kumo/collectors"
	sche "github.com/Khaym03/kumo/scheduler"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type Kumo interface {
	Run()
	RegisterCollector(c collectors.Collector)
	Shutdown()
}

type KumoBrowser struct {
	ctx       context.Context
	browser   *rod.Browser
	scheduler sche.Scheduler
	*collectors.CollectorRegistry
}

func NewKumoBrowser() *KumoBrowser {
	u := launcher.New().Leakless(false).MustLaunch()
	b := rod.New().ControlURL(u).Trace(true).MustConnect()

	return &KumoBrowser{
		ctx:               context.Background(),
		browser:           b,
		scheduler:         sche.NewScheduler(b, 2),
		CollectorRegistry: collectors.NewCollectorRegistry(),
	}
}

func (k *KumoBrowser) Run() {
	for _, c := range k.Collectors {
		err := c.Collect(k.ctx)
		if err != nil {
			log.Println(err)
		}
	}
}

func (k KumoBrowser) Shutdown() {
	k.browser.MustClose()
}
