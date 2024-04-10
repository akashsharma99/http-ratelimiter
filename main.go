package main

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

func main() {
	// create a new token bucket
	tb := &TokenBucket{
		capacity:   5,
		tokens:     0,
		rate:       1,
		lastUpdate: time.Now().Unix(),
	}
	// create a basic http server
	http.HandleFunc("GET /resource", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!\n"))
		slog.Info("Request received", "path", r.URL.Path, "method", r.Method, "remote_addr", r.RemoteAddr)
	})
	// create a ticker to add tokens to the bucket
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	go func() {
		for t := range ticker.C {
			tb.AddTokens(t)
		}
	}()
	http.ListenAndServe(":8080", nil)
}
