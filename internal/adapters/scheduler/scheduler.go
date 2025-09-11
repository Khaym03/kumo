package pagepool

import (
	"github.com/Khaym03/kumo/internal/pkg/browser"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// PageSetupFunc is a "hook" function type.
type PageSetupFunc func(page *rod.Page) error

// PagePool manages a pool of pages by using a browser pool.
type PagePool struct {
	browserPool  *browser.BrowserPool
	pool         rod.Pool[rod.Page]
	pageSetupFns []PageSetupFunc
}

// NewPagePool creates a new PagePool instance.
func NewPagePool(bp *browser.BrowserPool, numOfPagePerBrowser int, pageSetupFns ...PageSetupFunc) *PagePool {
	size := numOfPagePerBrowser * bp.Size()

	return &PagePool{
		browserPool:  bp,
		pool:         rod.NewPagePool(size),
		pageSetupFns: pageSetupFns,
	}
}

// Get gets a page from the pool, creating one if needed.
func (pp *PagePool) Get() (*rod.Page, error) {
	return pp.pool.Get(func() (*rod.Page, error) {
		browser, err := pp.browserPool.Get()
		if err != nil {
			return nil, err
		}
		defer pp.browserPool.Put(browser)

		page, err := browser.Page(proto.TargetCreateTarget{})
		if err != nil {
			return nil, err
		}

		for _, setupFunc := range pp.pageSetupFns {
			if err := setupFunc(page); err != nil {
				page.Close()
				return nil, err
			}
		}

		return page, nil
	})
}

// Put returns a page to the pool. The user is responsible for
// also returning the browser via page.Browser().
func (pp *PagePool) Put(p *rod.Page) {
	pp.pool.Put(p)
}
