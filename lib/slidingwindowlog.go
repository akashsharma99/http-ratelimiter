package lib

import (
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// going to use sliding window log to implement rate limiting
type SlidingWindow struct {
	// slice to store the timestamp of each request
	log []time.Time // will act as a queue
	// window size
	size int
	// mutex
	mu sync.Mutex
	// window duration
	duration time.Duration
}

func NewSlidingWindow(size int, duration time.Duration) *SlidingWindow {
	return &SlidingWindow{
		size:     size,
		duration: duration,
	}
}

func (sw *SlidingWindow) Allow() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	now := time.Now()
	// loop through the existing logs stamps and remove logs older than now - duration
	for len(sw.log) > 0 && now.Sub(sw.log[0]) > sw.duration {
		sw.log = sw.log[1:]
	}

	// check for window size
	if len(sw.log) >= sw.size {
		slog.Info("window size full :", "length", len(sw.log))
		return false
	}

	// add a new timestamp to the log
	sw.log = append(sw.log, now)
	slog.Info("added request log to window :", "timestamp", now.String())
	return true
}

func (sw *SlidingWindow) RateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if sw.Allow() {
			next(w, r)
		} else {
			http.Error(w, "Oooh slow down buddy!!", http.StatusTooManyRequests)
		}
	}
}
