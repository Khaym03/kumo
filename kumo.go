package main

import (
	sche "github.com/Khaym03/kumo/scheduler"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type Kumo struct {
	browser   *rod.Browser
	scheduler sche.Scheduler
	router    *rod.HijackRouter
}

func NewKumo() *Kumo {
	u := launcher.New().Leakless(false).MustLaunch()

	b := rod.New().ControlURL(u).Trace(true).MustConnect()

	return &Kumo{
		browser:   b,
		scheduler: *sche.NewScheduler(b, 2),
	}
}

func (k *Kumo) RegisterRouterMiddleware(pattern string, hander func(*rod.Hijack)) {
	k.router.MustAdd(pattern, hander)
}

func (k Kumo) Shutdown() {
	k.router.MustStop()
	k.browser.MustClose()
}
