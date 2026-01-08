package core

import "github.com/Khaym03/kumo/internal/ports"

type ResourcePool struct {
	Browsers ports.BrowserPool
	Pages    ports.PagePool
}

func NewResourcePool(bp ports.BrowserPool, pp ports.PagePool) *ResourcePool {
	return &ResourcePool{Browsers: bp, Pages: pp}
}

func (rp *ResourcePool) Close() error {
	return rp.Browsers.Close()
}