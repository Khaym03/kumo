package adapters

import (
	"context"
	"sync"
	"time"

	"github.com/Khaym03/kumo/internal/pkg/types"
	"github.com/Khaym03/kumo/internal/ports"
	"github.com/go-rod/rod"
	log "github.com/sirupsen/logrus"
)

type MockCollector struct{}

func NewMockCollector() ports.Collector {
	return new(MockCollector)
}

func (m MockCollector) String() string {
	return "Mock-collector"
}

func (m MockCollector) ProcessPage(
	ctx context.Context,
	page *rod.Page,
	req *types.Request,
	queue ports.Enqueuer,
	fs ports.FileStorage,
) error {
	log.Info("collecting something")
	page.MustNavigate(req.URL)

	time.Sleep(5 * time.Second)

	once.Do(func() {
		queue.Enqueue(&types.Request{URL: "https://proxyscrape.com/free-proxy-list", Collector: "Mock-collector"})
	})

	return nil
}

var once sync.Once
