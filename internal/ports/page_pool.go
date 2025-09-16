package ports

import (
	"github.com/Khaym03/kumo/internal/pkg/types"
	"github.com/go-rod/rod"
)

type PagePool interface {
	Get() (*rod.Page, error)
	Put(*rod.Page)
	Size() int
}

type BrowserPool interface {
	Get() (*rod.Browser, error)
	Put(*rod.Browser)
	Size() int
	Close() error
}

type PersistenceStore interface {
	SavePending(requests ...*types.Request) error
	LoadPending() ([]*types.Request, error)
	SaveCompleted(req *types.Request) error
	RemoveFromPending(req *types.Request) error
	// url or any other identifier
	IsCompleted(url string) (bool, error)

	Close() error
}

type FileStorage interface {
	SaveHTML(key string, data []byte) error
	SavePDF(key string, data []byte) error
	SaveJSON(id string, data []byte) error
	GetHTML(key string) ([]byte, error)
	GetPDF(key string) ([]byte, error)
	GetJSON(key string) ([]byte, error)
}
