package scheduler

import (
	"github.com/go-rod/rod"
)

type Scheduler struct {
	browser *rod.Browser
	rod.Pool[rod.Page]
}

func NewScheduler(b *rod.Browser, limit int) *Scheduler {
	pool := rod.NewPagePool(limit)

	return &Scheduler{
		browser: b,
		Pool:    pool,
	}
}
