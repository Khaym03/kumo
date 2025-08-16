package core

import (
	"context"
	"log"

	"github.com/Khaym03/kumo/collectors"
	sche "github.com/Khaym03/kumo/scheduler"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"

	_ "github.com/PuerkitoBio/goquery"
	_ "golang.org/x/time/rate"
)

type Kumo struct {
	ctx       context.Context
	browser   *rod.Browser
	scheduler sche.Scheduler
	*collectors.CollectorRegistry
}

func NewKumo() *Kumo {
	u := launcher.New().Leakless(false).MustLaunch()
	b := rod.New().ControlURL(u).Trace(true).MustConnect()

	return &Kumo{
		ctx:               context.Background(),
		browser:           b,
		scheduler:         sche.NewScheduler(b, 2),
		CollectorRegistry: collectors.NewCollectorRegistry(),
	}
}

func (k *Kumo) Run() {
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
