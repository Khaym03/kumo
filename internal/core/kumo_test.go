package core

import (
	"context"
	"errors"
	"testing"

	mocks "github.com/Khaym03/kumo/internal/mock/ports"
	"github.com/Khaym03/kumo/internal/pkg/types"
	"github.com/go-rod/rod"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestKumoEngine_ExecuteTask(t *testing.T) {
	mockDispatcher := mocks.NewMockDispatcher(t)
	mockPagePool := mocks.NewMockPagePool(t)
	mockCollector := mocks.NewMockCollector(t)
	
	mockPage := &rod.Page{}

	req := &types.Request{URL: "http://example.com", Collector: "test"}

	// We use .Maybe() because multiple calls might happen during init/test runs.
	mockCollector.On("String").Return("test").Maybe()
	mockPagePool.On("Size").Return(1).Maybe()

	resources := &ResourcePool{Pages: mockPagePool}
	engine := NewKumoEngine(context.Background(), resources, mockDispatcher, mockCollector)

	t.Run("successful processing", func(t *testing.T) {
		mockPagePool.On("Get").Return(mockPage, nil).Once()
		mockCollector.On("ProcessPage", mock.Anything, mockPage, req, mockDispatcher).Return(nil).Once()
		mockPagePool.On("Put", mockPage).Return().Once()
		mockDispatcher.On("Finish", req, nil).Return().Once()

		engine.executeTask(req)
	})

	t.Run("failed processing", func(t *testing.T) {
		errProcess := errors.New("scraping failed")

		mockPagePool.On("Get").Return(mockPage, nil).Once()
		mockCollector.On("ProcessPage", mock.Anything, mockPage, req, mockDispatcher).Return(errProcess).Once()
		mockPagePool.On("Put", mockPage).Return().Once()
		mockDispatcher.On("Finish", req, errProcess).Return().Once()

		engine.executeTask(req)
	})
}

func TestKumoEngine_Run(t *testing.T) {
	mockDispatcher := mocks.NewMockDispatcher(t)
	mockPagePool := mocks.NewMockPagePool(t)
	ctx := context.Background()
	
	// Expectations for constructor
	mockPagePool.On("Size").Return(1).Maybe()
	
	resources := &ResourcePool{Pages: mockPagePool}
	engine := NewKumoEngine(ctx, resources, mockDispatcher)

	t.Run("should load pending tasks and start workers", func(t *testing.T) {
		pending := []*types.Request{{URL: "pending-url"}}
		seeds := []*types.Request{{URL: "seed-url"}}
		
		mockDispatcher.On("LoadSavedState").Return(pending, nil).Once()
		
		// We expect the dispatch of combined requests
		mockDispatcher.On("Dispatch", mock.Anything, mock.Anything).Return().Once()

		mockDispatcher.On("Pull", mock.Anything).Return(nil, false).Maybe()

		// The engine calls Shutdown to close the queue
		mockDispatcher.On("Shutdown").Return().Once()

		err := engine.Run(seeds...)
		assert.NoError(t, err)
	})
}