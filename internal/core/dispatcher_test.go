package core

import (
	"errors"
	"testing"

	mocks "github.com/Khaym03/kumo/internal/mock/ports"
	"github.com/Khaym03/kumo/internal/pkg/types"
	"github.com/stretchr/testify/mock"
)

func TestDispatcher_Finish(t *testing.T) {
	mockFilter := mocks.NewMockRequestFilter(t)
	mockStore := mocks.NewMockPersistenceStore(t)
	mockLimiter := mocks.NewMockRateLimiter(t)
	dispatcher := NewDispatcher(mockFilter, mockStore, mockLimiter)

	req := &types.Request{URL: "http://test.com"}

	t.Run("should not save completed if error is present", func(t *testing.T) {
		// Mock internal WG manually or through behavior
		dispatcher.wg.Add(1)

		errTest := errors.New("some error")
		// We expect NO calls to SaveCompleted

		dispatcher.Finish(req, errTest)

		// If the test doesn't hang here, wg.Done() was called correctly
		dispatcher.wg.Wait()
		mockStore.AssertNotCalled(t, "SaveCompleted", mock.Anything)
	})
}