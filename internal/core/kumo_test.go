package core_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/Khaym03/kumo/internal/core"
	mockports "github.com/Khaym03/kumo/internal/mock/ports"
	"github.com/Khaym03/kumo/internal/pkg/types"
)

// KumoEngineTestSuite defines the suite for KumoEngine tests.
type KumoEngineTestSuite struct {
	suite.Suite
	ctx             context.Context
	cancel          context.CancelFunc
	engine          *core.KumoEngine
	mockReqStore    *mockports.MockPersistenceStore
	mockPagePool    *mockports.MockPagePool
	mockBrowserPool *mockports.MockBrowserPool
	mockCollector   *mockports.MockCollector
	mockFileStorage *mockports.MockFileStorage
}

// SetupTest is called before each test method in the suite.
func (s *KumoEngineTestSuite) SetupTest() {
	s.ctx, s.cancel = context.WithTimeout(context.Background(), 5*time.Second)

	s.mockReqStore = new(mockports.MockPersistenceStore)
	s.mockPagePool = new(mockports.MockPagePool)
	s.mockBrowserPool = new(mockports.MockBrowserPool)
	s.mockCollector = new(mockports.MockCollector)
	s.mockFileStorage = new(mockports.MockFileStorage)

	s.mockPagePool.On("Size").Return(1)
	s.mockCollector.On("String").Return("test-collector")
	s.mockReqStore.On("IsCompleted", mock.Anything).Return(false, nil)
	s.mockReqStore.On("SavePending", mock.Anything).Return(nil)
	s.mockReqStore.On("SaveCompleted", mock.Anything).Return(nil)
	s.mockReqStore.On("RemoveFromPending", mock.Anything).Return(nil)

	s.engine = core.NewKumoEngine(
		s.ctx,
		s.mockBrowserPool,
		s.mockPagePool,
		s.mockReqStore,
		s.mockFileStorage,
		s.mockCollector,
	)
}

// TearDownTest is called after each test method in the suite.
func (s *KumoEngineTestSuite) TearDownTest() {
	s.cancel()
}

// TestSuccessfulFlow is a test method for the successful crawl scenario.
func (s *KumoEngineTestSuite) TestSuccessfulFlow() {
	// Arrange
	initialReqs := []*types.Request{
		{URL: "http://example.com/page1", Collector: "test-collector"},
		{URL: "http://example.com/page2", Collector: "test-collector"},
	}

	s.mockReqStore.On("LoadPending").Return([]*types.Request{}, nil).Once()

	s.mockPagePool.On("Get").Return(new(rod.Page), nil).Twice()
	s.mockPagePool.On("Put", mock.Anything).Return(nil).Twice()

	s.mockCollector.On("ProcessPage", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Twice()

	// Act
	err := s.engine.Run(initialReqs...)

	// Assert
	assert.NoError(s.T(), err, "Run() should not return an error")
	s.mockReqStore.AssertExpectations(s.T())
	s.mockPagePool.AssertExpectations(s.T())
	s.mockCollector.AssertExpectations(s.T())
}

// TestResumesFromPending is a test method for the resume-from-pending scenario.
func (s *KumoEngineTestSuite) TestResumesFromPending() {
	// Arrange
	pendingReqs := []*types.Request{
		{URL: "http://example.com/pending", Collector: "test-collector"},
	}

	s.mockReqStore.On("LoadPending").Return(pendingReqs, nil).Once()

	s.mockPagePool.On("Get").Return(new(rod.Page), nil).Once()
	s.mockPagePool.On("Put", mock.Anything).Return(nil).Once()

	s.mockCollector.On("ProcessPage", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	// Act
	err := s.engine.Run()

	// Assert
	assert.NoError(s.T(), err, "Run() should not return an error")
	s.mockReqStore.AssertExpectations(s.T())
	s.mockPagePool.AssertExpectations(s.T())
	s.mockCollector.AssertExpectations(s.T())
}

// TestKumoEngine is the entry point for running the test suite.
func TestKumoEngine(t *testing.T) {
	suite.Run(t, new(KumoEngineTestSuite))
}
