package adapters

import (
	"context"
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

func (m MockCollector) Name() string {
	return "Mock-collector"
}

func (m MockCollector) ProcessPage(
	ctx context.Context,
	page *rod.Page,
	req *types.Request,
	queue ports.Enqueuer,
	ds ports.DataSink,
) error {
	log.Info("collecting something")
	page.MustNavigate(req.URL)

	time.Sleep(5 * time.Second)

	return nil
}
