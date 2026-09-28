package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	mu      sync.Mutex
	clients = map[string]*Client{}
)

type Client struct {
	tokens float64
	last   time.Time
}

func rateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		mu.Lock()

		c := clients[ip]
		if c == nil {
			c = &Client{5, time.Now()}
			clients[ip] = c
		}

		now := time.Now()
		c.tokens += now.Sub(c.last).Seconds() * 5
		if c.tokens > 5 {
			c.tokens = 5
		}
		c.last = now

		if c.tokens < 1 {
			mu.Unlock()
			http.Error(w, "Too Many Requests", 429)
			return
		}

		c.tokens--
		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

func main() {
	http.Handle("/", rateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Request accepted")
	})))

	fmt.Println("Server running on :8888")
	http.ListenAndServe(":8888", nil)
}