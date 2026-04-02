package network

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	global   *rate.Limiter
	qps      float64
}

func NewRateLimiter(qps float64) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		global:   rate.NewLimiter(rate.Limit(qps), 1),
		qps:      qps,
	}
}

func (r *RateLimiter) GetLimiter(key string) *rate.Limiter {
	r.mu.Lock()
	defer r.mu.Unlock()

	if limiter, exists := r.limiters[key]; exists {
		return limiter
	}

	limiter := rate.NewLimiter(rate.Limit(r.qps), 1)
	r.limiters[key] = limiter
	return limiter
}

func (r *RateLimiter) Wait(ctx context.Context, key string) error {
	limiter := r.GetLimiter(key)
	return limiter.Wait(ctx)
}

func (r *RateLimiter) Allow(key string) bool {
	limiter := r.GetLimiter(key)
	return limiter.Allow()
}

func (r *RateLimiter) SetQPS(qps float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.qps = qps
	r.global = rate.NewLimiter(rate.Limit(qps), 1)
}

func (r *RateLimiter) GetQPS() float64 {
	return r.qps
}

type AdaptiveLimiter struct {
	limiter    *rate.Limiter
	mu         sync.Mutex
	baseQPS    float64
	currentQPS float64
	minQPS     float64
	maxQPS     float64
	increaseAt time.Time
	decreaseAt time.Time
}

func NewAdaptiveLimiter(baseQPS float64) *AdaptiveLimiter {
	return &AdaptiveLimiter{
		limiter:    rate.NewLimiter(rate.Limit(baseQPS), 1),
		baseQPS:    baseQPS,
		currentQPS: baseQPS,
		minQPS:     baseQPS * 0.1,
		maxQPS:     baseQPS * 10,
		increaseAt: time.Now().Add(10 * time.Second),
		decreaseAt: time.Now().Add(60 * time.Second),
	}
}

func (a *AdaptiveLimiter) Wait(ctx context.Context) error {
	return a.limiter.Wait(ctx)
}

func (a *AdaptiveLimiter) Allow() bool {
	return a.limiter.Allow()
}

func (a *AdaptiveLimiter) Adjust(success bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := time.Now()

	if success {
		if now.After(a.increaseAt) && a.currentQPS < a.maxQPS {
			a.currentQPS *= 1.2
			if a.currentQPS > a.maxQPS {
				a.currentQPS = a.maxQPS
			}
			a.limiter.SetLimit(rate.Limit(a.currentQPS))
			a.increaseAt = now.Add(10 * time.Second)
		}
	} else {
		if now.After(a.decreaseAt) && a.currentQPS > a.minQPS {
			a.currentQPS *= 0.8
			if a.currentQPS < a.minQPS {
				a.currentQPS = a.minQPS
			}
			a.limiter.SetLimit(rate.Limit(a.currentQPS))
			a.decreaseAt = now.Add(5 * time.Second)
		}
	}
}

func (a *AdaptiveLimiter) GetCurrentQPS() float64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.currentQPS
}
