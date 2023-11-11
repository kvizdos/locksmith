package ratelimits

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	RequestsPerSecond    rate.Limit
	SecondsBurstCapacity int
	RequestsPerMinute    rate.Limit
	MinutesBurstCapacity int

	secondLimiters sync.Map // stores *rate.Limiter for per-second limits keyed by client ID
	minuteLimiters sync.Map // stores *rate.Limiter for per-minute limits keyed by client ID
}

// NewRateLimiter initializes a new RateLimiter with per-minute ratings
func NewRateLimiter(reqPerMin rate.Limit, minBurst int) *RateLimiter {
	return &RateLimiter{
		RequestsPerMinute:    reqPerMin,
		MinutesBurstCapacity: minBurst,
	}
}

func (rl *RateLimiter) WithSecondsLimits(reqPerSec rate.Limit, secBurst int) *RateLimiter {
	rl.RequestsPerSecond = reqPerSec
	rl.SecondsBurstCapacity = secBurst

	return rl
}

// getSecondLimiter returns the per-second rate limiter for the provided ID, creating a new one if necessary.
func (rl *RateLimiter) getSecondLimiter(id string) *rate.Limiter {
	limiter, exists := rl.secondLimiters.Load(id)
	if !exists {
		limiter = rate.NewLimiter(rl.RequestsPerSecond, rl.SecondsBurstCapacity)
		rl.secondLimiters.Store(id, limiter)
	}
	return limiter.(*rate.Limiter)
}

// getMinuteLimiter returns the per-minute rate limiter for the provided ID, creating a new one if necessary.
func (rl *RateLimiter) getMinuteLimiter(id string) *rate.Limiter {
	limiter, exists := rl.minuteLimiters.Load(id)
	if !exists {
		limiter = rate.NewLimiter(rate.Every(time.Minute/time.Duration(rl.RequestsPerMinute)), rl.MinutesBurstCapacity)
		rl.minuteLimiters.Store(id, limiter)
	}
	return limiter.(*rate.Limiter)
}

// CanHandleRequest checks if the rate limiters allow another request based on both per-second and per-minute limits.
func (rl *RateLimiter) CanHandle(identifier string) bool {
	// Check the per-second limiter first
	if rl.RequestsPerSecond > 0 {
		if !rl.getSecondLimiter(identifier).Allow() {
			fmt.Println("Per second not allowed!")
			return false
		}
	}
	// Check the per-minute limiter
	min := rl.getMinuteLimiter(identifier).Allow()

	return min
}
