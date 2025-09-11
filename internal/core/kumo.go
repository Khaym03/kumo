package core

import (
	"context"
	"log"

	"github.com/Khaym03/kumo/internal/ports"

	sche "github.com/Khaym03/kumo/scheduler"
	"github.com/go-rod/rod"

	_ "github.com/PuerkitoBio/goquery"
	_ "golang.org/x/time/rate"
)

type Kumo struct {
	ctx       context.Context
	browser   *rod.Browser
	scheduler sche.Scheduler
	ports.CollectorRegistry
}

func NewKumo(
	browser *rod.Browser,
	scheduler sche.Scheduler,
	registry ports.CollectorRegistry,
) *Kumo {
	return &Kumo{
		browser:           browser,
		scheduler:         scheduler,
		CollectorRegistry: registry,
	}
}

func (k *Kumo) Run() {
	for _, c := range k.Collectors() {
		err := c.Collect(k.ctx)
		if err != nil {
			log.Println(err)
		}
	}
}

func (k Kumo) Shutdown() {
	k.browser.MustClose()
}
