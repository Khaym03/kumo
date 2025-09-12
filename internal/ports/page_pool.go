package ports

import "github.com/go-rod/rod"

type PagePool interface {
	Get() (*rod.Page, error)
	Put(*rod.Page)
}

type BrowserPool interface {
	Get() (*rod.Browser, error)
	Put(*rod.Browser)
	Size() int
	Close() error
}
