package gre

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

/**
 * Package name: gre
 * Project name: go-routing-engine
 * Created by: Praveen Premaratne
 * Created on: 27/04/2023 21:23
 */

// NewRouter returns a mux.Router with given routes added
// pluggable to a standard http.Server handler
//
// param: <routes> a Routes object with list of Route objects
//
// param: <strictSlashes> defines the trailing slash behavior for new routes
func NewRouter(routes Routes, strictSlashes bool) *mux.Router {
	return addRoutes(routes, strictSlashes)
}

func addRoutes(routes Routes, strict bool) *mux.Router {
	router := mux.NewRouter().StrictSlash(strict)

	log.Println("add global handler 404 - not found")
	router.NotFoundHandler = http.HandlerFunc(add404)

	log.Println("add global handler 405 - method not allowed")
	router.MethodNotAllowedHandler = http.HandlerFunc(add405)

	for _, route := range routes {
		var handler http.Handler
		if route.Deprecated {
			log.Printf("ignore mapping: %s ( %s %s ) deprecated\n", route.Name, route.Methods, route.Pattern)
			handler = http.HandlerFunc(deprecated)
		} else {
			handler = route.HandlerFunc
		}

		handler = Logger(handler, route.Name)
		router.
			Methods(route.Methods...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

		log.Printf("add mapping: %s ( %s %s )\n", route.Name, route.Methods, route.Pattern)
	}

	router.
		Methods(http.MethodGet).
		Path("/metrics").
		Handler(promhttp.Handler())

	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(promMiddleware)

	return router
}

func add404(w http.ResponseWriter, r *http.Request) {
	resp := &ErrorResponse{
		Code:  http.StatusNotFound,
		Cause: "resource not found",
	}

	w.WriteHeader(resp.Code)
	fmt.Fprint(w, resp.Json())
}

func add405(w http.ResponseWriter, r *http.Request) {
	resp := &ErrorResponse{
		Code:  http.StatusMethodNotAllowed,
		Cause: "method not allowed",
	}

	w.WriteHeader(resp.Code)
	fmt.Fprint(w, resp.Json())
}

func deprecated(w http.ResponseWriter, r *http.Request) {
	resp := &ErrorResponse{
		Code:  http.StatusForbidden,
		Cause: "method deprecated",
	}

	w.WriteHeader(resp.Code)
	fmt.Fprint(w, resp.Json())
}
