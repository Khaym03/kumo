package scheduler

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

type Scheduler interface {
	Get() (*rod.Page, error)
	Put(*rod.Page)
}

type scheduler struct {
	browser *rod.Browser
	pool    rod.Pool[rod.Page]
}

func NewScheduler(b *rod.Browser, limit int) *scheduler {
	pool := rod.NewPagePool(limit)

	return &scheduler{
		browser: b,
		pool:    pool,
	}
}

func (s scheduler) Get() (*rod.Page, error) {
	return s.pool.Get(func() (*rod.Page, error) {
		page, err := s.browser.Page(proto.TargetCreateTarget{})
		if err != nil {
			return nil, err
		}

		_, err = page.EvalOnNewDocument(stealth.JS)
		if err != nil {
			return nil, err
		}

		return page, nil
	})
}

func (s scheduler) Put(p *rod.Page) {
	s.pool.Put(p)
}
