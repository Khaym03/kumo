package main

import (
	adapters "github.com/Khaym03/kumo/internal/adapters/collector"
	"github.com/Khaym03/kumo/internal/adapters/config"
	"github.com/Khaym03/kumo/internal/adapters/pagepool"
	sinks "github.com/Khaym03/kumo/internal/adapters/storage"
	"github.com/Khaym03/kumo/internal/core"
	"github.com/Khaym03/kumo/internal/pkg/browser"
	"github.com/Khaym03/kumo/internal/pkg/proxy"
	"github.com/Khaym03/kumo/internal/pkg/types"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("loading config...")
	conf := config.LoadKumoConfig()

	log.Println(*conf)
	pm := proxy.NewConcurrentProxyManager([]proxy.Proxy{})

	creators := browser.CreateCreatorsFromConfig(conf.Browsers, pm)
	browserPool := browser.NewPool(creators...)
	pp := pagepool.NewPagePool(browserPool, conf.NumOfPagesPerBrowser)

	ds, err := sinks.NewHTMLFileSink("test/")
	if err != nil {
		log.Fatal(err)
	}

	mockCollector := adapters.NewMockCollector()

	kumo := core.NewKumoEngine(browserPool, pp, ds, mockCollector)

	initialRequest := &types.Request{
		URL:       "https://example.com",
		Collector: mockCollector.Name(),
	}

	kumo.InitialReqs(initialRequest)

	kumo.Run()
	// kumo.Shutdown()

}
