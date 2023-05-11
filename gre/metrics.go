package gre

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
)

/**
 * Package name: gre
 * Project name: go-routing-engine
 * Created by: Praveen Premaratne
 * Created on: 28/04/2023 23:06
 */

var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of get requests.",
		},
		[]string{"path"},
	)

	httpDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_response_time_seconds",
			Help: "Duration of HTTP requests.",
		}, []string{"path"},
	)
)

func promMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(r.URL.Path))
		totalRequests.WithLabelValues(r.URL.Path).Inc()

		next.ServeHTTP(w, r)

		timer.ObserveDuration()
	})
}
