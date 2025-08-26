package composer

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Khaym03/kumo/config"
	"github.com/Khaym03/kumo/controller"
	"github.com/Khaym03/kumo/proxy"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
)

var ErrOutOfProxies = errors.New("no available proxies to use")

type proxiedConnection struct {
	proxy.Proxy
	isInUse bool
}

func newProxiedConnection(p proxy.Proxy) *proxiedConnection {
	return &proxiedConnection{Proxy: p}
}

type BrowserFactory struct {
	*config.RemoteConfig
	rod.Pool[rod.Browser]
	size    int
	proxies []*proxiedConnection
	mu      sync.Mutex
}

func NewBroserFactory(
	c *config.RemoteConfig,
	proxies []proxy.Proxy,
	size int,
) *BrowserFactory {

	pcs := make([]*proxiedConnection, len(proxies))
	for i, prox := range proxies {
		pcs[i] = newProxiedConnection(prox)
	}

	return &BrowserFactory{
		RemoteConfig: c,
		proxies:      pcs,
		size:         size,
		Pool:         rod.NewBrowserPool(size),
		mu:           sync.Mutex{},
	}
}

func (bf *BrowserFactory) LocalBrowserCreator() (*rod.Browser, error) {
	// Leckless is set to false since is detected as threat in windows
	// so must handle the close
	u, err := launcher.New().Leakless(false).Launch()
	if err != nil {
		return nil, err
	}

	browser := rod.New().ControlURL(u).Trace(true)
	err = browser.Connect()
	if err != nil {
		return nil, err
	}

	return browser, nil
}

func (bf *BrowserFactory) LocalBrowserWithProxyCreator() (*rod.Browser, error) {
	launcherFunc := func(p *proxiedConnection) (string, error) {
		l := launcher.New().Leakless(false)
		l.Set(flags.ProxyServer, p.Address())
		return l.Launch()
	}
	return bf.createBrowserWithProxy(launcherFunc)
}

func (bf *BrowserFactory) RemoteBrowserCreator() (*rod.Browser, error) {
	launcherFunc := func(p *proxiedConnection) (string, error) {
		machineAddress := fmt.Sprintf("%s:%d", bf.RemoteHosts[0], bf.RemotePort)
		return controller.NewWSURLBuilder(machineAddress).WithProxy(p.Proxy).Build()
	}

	return bf.createBrowserWithProxy(launcherFunc)
}

func (bf *BrowserFactory) createBrowserWithProxy(
	launcherFunc func(p *proxiedConnection) (string, error),
) (*rod.Browser, error) {
	p := bf.getAnAvailableProxy()
	if p == nil {
		return nil, ErrOutOfProxies
	}

	wsURL, err := launcherFunc(p)
	if err != nil {
		bf.ReleaseProxy(p)
		return nil, err
	}

	browser := rod.New().ControlURL(wsURL).Trace(true)
	err = browser.Connect()
	if err != nil {
		bf.ReleaseProxy(p)
		return nil, err
	}

	go browser.MustHandleAuth(p.User, p.Password)()

	browser.MustIgnoreCertErrors(true)

	return browser, nil
}

func (bf *BrowserFactory) getAnAvailableProxy() *proxiedConnection {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	for _, p := range bf.proxies {
		if !p.isInUse {
			p.isInUse = true
			return p
		}
	}
	return nil
}

func (bf *BrowserFactory) ReleaseProxy(p *proxiedConnection) {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	p.isInUse = false
}
