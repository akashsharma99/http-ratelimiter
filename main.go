package main

import (
	"log/slog"
	"net/http"
	"time"
)

type TokenBucket struct {
	// The total number of tokens the bucket can hold.
	Capacity int
	// The current number of tokens in the bucket.
	Tokens int
	// Number of tokens to add per second.
	Rate int
	// Last time the bucket was updated.
	LastUpdate int64
}

func (tb *TokenBucket) AddTokens() {
	now := time.Now().Unix()
	tokensToAdd := int(now-tb.LastUpdate) * tb.Rate
	tb.Tokens = tb.Tokens + tokensToAdd
	if tb.Tokens > tb.Capacity {
		tb.Tokens = tb.Capacity
	}
	tb.LastUpdate = now
}

func main() {
	// create a new token bucket
	tb := &TokenBucket{
		Capacity:   5,
		Tokens:     0,
		Rate:       1,
		LastUpdate: time.Now().Unix(),
	}
	// create a basic http server
	http.HandleFunc("GET /resource", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!\n"))
		slog.Info("Request received", "path", r.URL.Path, "method", r.Method, "remote_addr", r.RemoteAddr)
	})
	http.ListenAndServe(":8080", nil)
}
