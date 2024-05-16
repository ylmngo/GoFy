package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Middleware func(next http.Handler) http.Handler

func (app *application) recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				er := fmt.Sprintf("%s", err)
				app.writeJSON(w, http.StatusInternalServerError, er, nil)
				app.logger.Printf("Internal Server Error\n")
				w.Header().Set("Connection", "close")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.cfg.limiter.enabled {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.writeJSON(w, http.StatusInternalServerError, "server encountered a problem", nil)
				return
			}

			mu.Lock()

			if _, ok := clients[ip]; !ok {
				clients[ip] = &client{
					limiter:  rate.NewLimiter(2, 4),
					lastSeen: time.Now(),
				}
			}

			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.writeJSON(w, http.StatusTooManyRequests, "rate limit exceeded", nil)
				return
			}

			mu.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}
