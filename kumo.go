package main

import (
	"context"
	"log"

	"github.com/Khaym03/kumo/collectors"
	sche "github.com/Khaym03/kumo/scheduler"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type Kumo struct {
	ctx        context.Context
	browser    *rod.Browser
	scheduler  sche.Scheduler
	collectors []collectors.Collector
}

func NewKumo() *Kumo {
	u := launcher.New().Leakless(false).MustLaunch()
	b := rod.New().ControlURL(u).Trace(true).MustConnect()

	return &Kumo{
		ctx:        context.Background(),
		browser:    b,
		scheduler:  sche.NewScheduler(b, 2),
		collectors: []collectors.Collector{},
	}
}

func (k *Kumo) Run() {
	for _, c := range k.collectors {
		err := c.Collect(k.ctx)
		if err != nil {
			log.Println(err)
		}
	}
}

func (k *Kumo) RegisterCollector(c collectors.Collector) {
	k.collectors = append(k.collectors, c)
}

func (k Kumo) Shutdown() {
	k.browser.MustClose()
}
