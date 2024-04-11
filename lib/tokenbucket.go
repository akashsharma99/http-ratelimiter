package lib

import (
	"log/slog"
	"net/http"
	"time"
)

type TokenBucket struct {
	// The total number of tokens the bucket can hold.
	capacity int
	// The current number of tokens in the bucket.
	tokens int
	// Number of tokens to add per second.
	rate int
	// Last time the bucket was updated.
	lastUpdate int64
}

func NewTokenBucket(capacity int, rate int, initialTokens int) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     initialTokens,
		rate:       rate,
		lastUpdate: time.Now().Unix(),
	}
}

func (tb *TokenBucket) AddTokens(t time.Time) {
	// if bucket already full then no need to add tokens
	if tb.tokens == tb.capacity {
		return
	}
	tb.tokens += int(t.Unix()-tb.lastUpdate) * tb.rate
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
	tb.lastUpdate = t.Unix()
	slog.Info("Tokens added", "tokens", tb.tokens, "last_update", tb.lastUpdate)
}

// middleware func to take token from bucket
func (tb *TokenBucket) RateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if tb.tokens == 0 {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		tb.tokens--
		next(w, r)
	}
}
