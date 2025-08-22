package collectors

import "context"

type MockCollector struct{}

func (mc *MockCollector) Collect(ctx context.Context) error {
	return nil
}
