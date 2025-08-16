package main

import (
	"context"

	"github.com/Khaym03/kumo/core"
)

type MockCollector struct {
}

func (mc *MockCollector) Collect(ctx context.Context) error {
	return nil
}

func main() {
	kumo := core.NewKumo()
	defer kumo.Shutdown()

	kumo.RegisterCollector(&MockCollector{})

	kumo.Run()
}
