package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/akashsharma99/http-ratelimiter/lib"
)

func main() {
	//tryTokenBucket()
	tryFixedWindowCounter()
}
func tryTokenBucket() {
	// create a new token bucket
	tb := lib.NewTokenBucket(10, 1, 10)
	// create a basic http server
	http.HandleFunc("GET /resource", tb.RateLimit(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!\n"))
		slog.Info("Request received", "path", r.URL.Path, "method", r.Method, "remote_addr", r.RemoteAddr)
	}))
	// create a ticker to add tokens to the bucket
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	go func() {
		for t := range ticker.C {
			tb.AddTokens(t)
		}
	}()
	slog.Info("Server started", "port", 8080)
	http.ListenAndServe(":8080", nil)
}

func tryFixedWindowCounter() {
	// create a new fixed window counter
	fwc := lib.NewFixedWindowCounter(60, 60)
	// create a basic http server
	http.HandleFunc("GET /resource", fwc.RateLimit(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!\n"))
		slog.Info("Request received", "path", r.URL.Path, "method", r.Method, "remote_addr", r.RemoteAddr)
	}))
	slog.Info("Server started", "port", 8080)
	http.ListenAndServe(":8080", nil)
}
