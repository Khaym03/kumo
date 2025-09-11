package browser

import (
	"log"
	"sync"

	"github.com/Khaym03/kumo/internal/adapters/config"
	"github.com/Khaym03/kumo/internal/pkg/proxy"

	"github.com/go-rod/rod"
)

// BrowserFactory orchestrates the creation and pooling of browsers.
type BrowserPool struct {
	pool     rod.Pool[rod.Browser]
	creators []BrowserCreator
	mu       sync.Mutex
	next     int
}

// NewBrowserFactory creates a BrowserFactory. It infers the pool size
// from the number of BrowserCreator instances provided.
func NewPool(creators ...BrowserCreator) *BrowserPool {
	if len(creators) == 0 {
		log.Fatal("No browser creators provided. Cannot create a browser factory.")
	}

	return &BrowserPool{
		creators: creators,
		pool:     rod.NewPool[rod.Browser](len(creators)),
		mu:       sync.Mutex{},
	}
}

// Get gets a browser from the pool, creating one if needed.
func (bp *BrowserPool) Get() (*rod.Browser, error) {
	return bp.pool.Get(func() (*rod.Browser, error) {
		// Select a creator in a thread-safe, round-robin manner.
		bp.mu.Lock()
		creator := bp.creators[bp.next]
		bp.next = (bp.next + 1) % len(bp.creators)
		bp.mu.Unlock()

		return creator.CreateBrowser()
	})
}

func (bp *BrowserPool) Put(b *rod.Browser) {
	bp.pool.Put(b)
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
