package main

import (
	"context"
	"expvar"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/pascaldekloe/jwt"
	"golang.org/x/time/rate"
)

type Middleware func(next http.Handler) http.Handler
type contextKey string

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

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		token := r.Header.Get("Authorization")

		if token == "" {
			app.writeJSON(w, http.StatusNotAcceptable, "authorization header must be provided", nil)
			return
		}

		claims, err := jwt.HMACCheck([]byte(token), []byte(app.cfg.jwtSec))
		if err != nil {
			app.writeJSON(w, http.StatusBadRequest, "invalid authentication token", nil)
			return
		}

		if !claims.Valid(time.Now()) {
			app.writeJSON(w, http.StatusBadRequest, "invalid authentication token", nil)
			return
		}

		if claims.Issuer != "http://localhost:8000" {
			app.writeJSON(w, http.StatusBadRequest, "invalid issuer", nil)
			return
		}

		id, _ := strconv.ParseInt(claims.Subject, 10, 64)

		user, err := app.models.Users.GetById(int(id))
		if err != nil {
			app.writeJSON(w, http.StatusNotFound, "Invalid User", nil)
			app.logger.Printf("Unable to get user by id: %v\n", err)
			return
		}

		ctx := context.WithValue(r.Context(), contextKey("userID"), user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (app *application) metrics(next http.Handler) http.Handler {
	requests := expvar.NewInt("total_requests_recieved")
	responses := expvar.NewInt("total_responses_sent")
	totalProcessingTime := expvar.NewInt("total_process_time")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requests.Add(1)
		next.ServeHTTP(w, r)
		responses.Add(1)
		d := time.Since(start).Milliseconds()

		app.logger.Printf("%v --- %v ms\n", r.URL.Path, d)

		totalProcessingTime.Add(d)
	})
}
