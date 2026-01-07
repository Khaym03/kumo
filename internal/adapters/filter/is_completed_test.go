package filter

import (
	"errors"
	"testing"

	mocks "github.com/Khaym03/kumo/internal/mock/ports"
	"github.com/Khaym03/kumo/internal/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIsCompletedFilter_Filter(t *testing.T) {
	// Define the test cases
	tests := []struct {
		name           string
		inputReqs      []*types.Request
		setupMock      func(m *mocks.MockCompletionChecker)
		expectedLength int
	}{
		{
			name: "should filter out completed requests",
			inputReqs: []*types.Request{
				{URL: "https://done.com"},
				{URL: "https://pending.com"},
			},
			setupMock: func(m *mocks.MockCompletionChecker) {
				m.On("IsCompleted", "https://done.com").Return(true, nil)
				m.On("IsCompleted", "https://pending.com").Return(false, nil)
			},
			expectedLength: 1, // Only pending.com should remain
		},
		{
			name: "should skip request when checker returns error",
			inputReqs: []*types.Request{
				{URL: "https://error.com"},
			},
			setupMock: func(m *mocks.MockCompletionChecker) {
				m.On("IsCompleted", "https://error.com").Return(false, errors.New("db error"))
			},
			expectedLength: 0, // Code continues and doesn't append on error
		},
		{
			name: "should keep all when none are completed",
			inputReqs: []*types.Request{
				{URL: "https://a.com"},
				{URL: "https://b.com"},
			},
			setupMock: func(m *mocks.MockCompletionChecker) {
				m.On("IsCompleted", mock.Anything).Return(false, nil)
			},
			expectedLength: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockChecker := new(mocks.MockCompletionChecker)
			tt.setupMock(mockChecker)
			f := NewIsCompletedFilter(mockChecker)

			// Act
			result := f.Filter(tt.inputReqs)

			// Assert
			assert.Len(t, result, tt.expectedLength)
			mockChecker.AssertExpectations(t)
		})
	}
}
