package proxy

import (
	"errors"
	"sync"
)

type ProxyManager interface {
	Get() (Proxy, error)
	Release(p Proxy)
}

var ErrOutOfProxies = errors.New("no available proxies to use")

type ConcurrentProxyManager struct {
	mu      *sync.Mutex
	proxies map[Proxy]bool
}

func NewConcurrentProxyManager(proxies []Proxy) *ConcurrentProxyManager {
	m := map[Proxy]bool{}
	for _, p := range proxies {
		m[p] = false
	}

	return &ConcurrentProxyManager{
		mu:      new(sync.Mutex),
		proxies: m,
	}
}

// Get an available proxy
func (c *ConcurrentProxyManager) Get() (Proxy, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for p, isInUse := range c.proxies {
		if !isInUse {
			c.proxies[p] = true
			return p, nil
		}
	}
	return Proxy{}, ErrOutOfProxies
}

func (c *ConcurrentProxyManager) Release(p Proxy) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.proxies[p] = false
}
