package core

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDomainLimiter_NoBurstBehavior(t *testing.T) {
	ctx := context.Background()
	// Set 100ms interval between requests
	interval := 100 * time.Millisecond
	
	t.Run("should strictly enforce interval when burst is 1", func(t *testing.T) {
		// Burst = 1 means no extra tokens are stored. 
		// Every request after the first MUST wait for the interval.
		dl := NewDomainLimiter(interval, 1)
		url := "https://example.com"

		// 1. First call is usually near-instant
		_ = dl.Wait(ctx, url)
		
		start := time.Now()
		
		// 2. Second call must wait at least 'interval'
		err := dl.Wait(ctx, url)
		assert.NoError(t, err)
		
		elapsed := time.Since(start)
		
		// If burst were 0 or >1, the behavior would change. 
		// With 1, it acts as a strict throttler.
		assert.GreaterOrEqual(t, elapsed.Milliseconds(), interval.Milliseconds(), 
            "The second request should have been delayed by the interval")
	})

	t.Run("should allow immediate parallel execution for different domains", func(t *testing.T) {
		dl := NewDomainLimiter(1 * time.Second, 1)
		
		_ = dl.Wait(ctx, "https://a.com")
		
		start := time.Now()
		// This should NOT wait for a.com's second since it's a different bucket
		err := dl.Wait(ctx, "https://b.com")
		assert.NoError(t, err)
		
		assert.Less(t, time.Since(start), 50*time.Millisecond, 
            "Different domains should have independent buckets")
	})
}