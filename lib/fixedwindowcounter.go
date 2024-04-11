package lib

import (
	"net/http"
	"time"
)

// implementing fixed window counter rate limiter

type FixedWindowCounter struct {
	// window size in seconds
	windowSize int64
	// max number of requests allowed in window
	maxRequests int
	// current window start time
	windowStart int64
	// current window request count
	requestCount int
}

func NewFixedWindowCounter(windowSize int, maxRequests int) *FixedWindowCounter {
	return &FixedWindowCounter{
		windowSize:   int64(windowSize),
		maxRequests:  maxRequests,
		windowStart:  time.Now().Unix(),
		requestCount: 0,
	}
}

func (fwc *FixedWindowCounter) Allow() bool {
	now := time.Now().Unix()
	// current time is within the current window
	if now < fwc.windowStart+fwc.windowSize {
		// if request count is less than max requests then allow
		if fwc.requestCount < fwc.maxRequests {
			fwc.requestCount++
			return true
		}
		// if request count is greater than max requests then deny the request
		return false
	}
	// current time is outside the current window
	fwc.windowStart = now
	fwc.requestCount = 1
	return true
}
func (fwc *FixedWindowCounter) RateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if fwc.Allow() {
			next(w, r)
		} else {
			http.Error(w, "Oooh slow down buddy!!", http.StatusTooManyRequests)
		}
	}
}
