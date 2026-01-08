package core

import (
	"context"
	"net/url"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type DomainLimiter struct {
	limiters     map[string]*rate.Limiter
	mu           sync.RWMutex
	interval     time.Duration
	defaultBurst int
}


func NewDomainLimiter(interval time.Duration, burst int) *DomainLimiter {
	return &DomainLimiter{
		limiters:     make(map[string]*rate.Limiter),
		interval:     interval,
		defaultBurst: burst,
	}
}

func (dl *DomainLimiter) Wait(ctx context.Context, rawURL string) error {
	domain := dl.getDomain(rawURL)
	return dl.getLimiter(domain).Wait(ctx)
}

func (dl *DomainLimiter) getDomain(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host == "" {
		return "default"
	}
	return parsed.Host
}

func (dl *DomainLimiter) getLimiter(domain string) *rate.Limiter {
	dl.mu.RLock()
	l, exists := dl.limiters[domain]
	dl.mu.RUnlock()

	if exists {
		return l
	}

	dl.mu.Lock()
	defer dl.mu.Unlock()

	if l, exists = dl.limiters[domain]; exists {
		return l
	}

	newLimiter := rate.NewLimiter(rate.Every(dl.interval), dl.defaultBurst)
	dl.limiters[domain] = newLimiter
	return newLimiter
}