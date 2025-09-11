package browser

import (
	"github.com/Khaym03/kumo/internal/adapters/remote"
	"github.com/Khaym03/kumo/internal/pkg/proxy"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
)

// BrowserCreator defines the contract for any browser creation method.
type BrowserCreator interface {
	CreateBrowser() (*rod.Browser, error)
}

// Option is a function that configures a browser creator.
type Option func(*Options)

// Options holds the configuration for a browser creation.
type Options struct {
	Proxy          *proxy.Proxy
	MachineAddress string
	Headless       bool
}

// WithProxy is a functional option to set a proxy address.
func WithProxy(p proxy.Proxy) Option {
	return func(o *Options) {
		o.Proxy = &p
	}
}

// WithRemoteHost is a functional option to set a remote host address.
func WithRemoteHost(address string) Option {
	return func(o *Options) {
		o.MachineAddress = address
	}
}

func WithHeadless(b bool) Option {
	return func(o *Options) {
		o.Headless = b
	}
}

// LocalBrowserCreator creates a local browser.
type LocalBrowserCreator struct {
	options *Options
}

// NewLocalBrowserCreator creates a new LocalBrowserCreator with functional options.
func NewLocalBrowserCreator(opts ...Option) *LocalBrowserCreator {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}
	return &LocalBrowserCreator{options: options}
}

func (l *LocalBrowserCreator) CreateBrowser() (*rod.Browser, error) {
	launcher := launcher.New().Leakless(false)
	if l.options.Proxy != nil {
		launcher.Set(flags.ProxyServer, l.options.Proxy.Address())
	}

	launcher.Headless(l.options.Headless)

	u, err := launcher.Launch()
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

// RemoteBrowserCreator creates a remote browser.
type RemoteBrowserCreator struct {
	options *Options
}

// NewRemoteBrowserCreator creates a new RemoteBrowserCreator with functional options.
func NewRemoteBrowserCreator(opts ...Option) *RemoteBrowserCreator {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}
	return &RemoteBrowserCreator{options: options}
}

func (r *RemoteBrowserCreator) CreateBrowser() (*rod.Browser, error) {
	var url string
	var err error

	if r.options.Proxy != nil {
		url, err = remote.NewWSURLBuilder(r.options.MachineAddress).
			WithProxy(*r.options.Proxy).
			Build()
	} else {
		url, err = remote.NewWSURLBuilder(r.options.MachineAddress).Build()
	}

	if err != nil {
		return nil, err
	}
	browser := rod.New().ControlURL(url).Trace(true)
	err = browser.Connect()
	if err != nil {
		return nil, err
	}

	return browser, nil
}
