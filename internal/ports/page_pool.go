package ports

import "github.com/go-rod/rod"

type PagePool interface {
	Get() (*rod.Page, error)
	Put(*rod.Page)
}
