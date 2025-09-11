package browser

import (
	"log"
	"sync"

	"github.com/Khaym03/kumo/config"
	"github.com/Khaym03/kumo/proxy"
	"github.com/go-rod/rod"
)

// BrowserFactory orchestrates the creation and pooling of browsers.
type BrowserFactory struct {
	rod.Pool[rod.Browser]
	creators []BrowserCreator
	mu       sync.Mutex
	next     int
}

// NewBrowserFactory creates a BrowserFactory. It infers the pool size
// from the number of BrowserCreator instances provided.
func NewFactory(creators ...BrowserCreator) *BrowserFactory {
	if len(creators) == 0 {
		log.Fatal("No browser creators provided. Cannot create a browser factory.")
	}

	return &BrowserFactory{
		creators: creators,
		Pool:     rod.NewBrowserPool(len(creators)),
		mu:       sync.Mutex{},
	}
}

// Get gets a browser from the pool, creating one if needed.
func (bf *BrowserFactory) Get() (*rod.Browser, error) {
	return bf.Pool.Get(func() (*rod.Browser, error) {
		// Select a creator in a thread-safe, round-robin manner.
		bf.mu.Lock()
		creator := bf.creators[bf.next]
		bf.next = (bf.next + 1) % len(bf.creators)
		bf.mu.Unlock()

		return creator.CreateBrowser()
	})
}

func CreateCreatorsFromConfig(configs []config.BrowserConfig, pm proxy.ProxyManager) []BrowserCreator {
	creators := make([]BrowserCreator, 0, len(configs))

	for _, bc := range configs {
		var creator BrowserCreator
		var opts []Option

		// If the config specifies a proxy, get one and add it as an option.
		if bc.Proxy {
			p, err := pm.Get()
			if err != nil {
				log.Printf("Failed to get proxy for browser %s: %v. Skipping creator.", bc.Name, err)
				continue
			}
			opts = append(opts, WithProxy(p))
		}

		// Create the correct type of creator and pass the options.
		switch bc.Type {
		case "local":
			creator = NewLocalBrowserCreator(opts...)
		case "remote":
			opts = append(opts, WithRemoteHost(bc.Host))
			creator = NewRemoteBrowserCreator(opts...)
		default:
			log.Fatalf("Unknown browser type: %s", bc.Type)
		}

		creators = append(creators, creator)
	}

	return creators
}
