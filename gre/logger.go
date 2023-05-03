package gre

import (
	"log"
	"net/http"
	"time"
)

/**
 * Package name: gre
 * Project name: go-routing-engine
 * Created by: Praveen Premaratne
 * Created on: 27/04/2023 23:25
 */

// Logger middleware will log all incoming request and the function that handled that request
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf(
			"%s %s %s %s %s %s",
			r.Header.Get("X-Real-IP"),
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
			r.UserAgent(),
		)

		inner.ServeHTTP(w, r)
	})
}
