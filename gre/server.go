package gre

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"time"
)

/**
 * Package name: gre
 * Project name: go-routing-engine
 * Created by: Praveen Premaratne
 * Created on: 28/04/2023 17:09
 */

var (
	// RouteTable can be used for storing Routes.
	//
	// if the routes are scattered on multiple files
	// this allows users to append Route to from anywhere
	// keeping the Route configurations and handler functions
	// on the same file
	RouteTable = Routes{
		Route{
			Name:        "Health",
			Methods:     []string{http.MethodGet},
			Pattern:     "/health",
			HandlerFunc: health,
		},
	}
)

// NewServer returns a vanilla Server definition for later
// configuration
func NewServer() *Server {
	return &Server{}
}

// DefaultServer is a pre-defined Server that can be used to
// quickly initialise an HTTP server with basic configuration
//
// param: <port> server port to listen
//
// param: <strictSlashes> boolean value for strict slashes
func DefaultServer(port int, strictSlashes bool) *Server {
	var server Server

	server.Addr = fmt.Sprintf("0.0.0.0:%d", port)
	server.ReadTimeout = time.Duration(15) * time.Second
	server.WriteTimeout = time.Duration(15) * time.Second
	server.StrictSlash = strictSlashes

	return &server
}

// AddMiddleware is for adding middleware to the router for
// customisation such as authentication and validation
//
// param: <middleware> is http.Handler method
func (s *Server) AddMiddleware(middleware func(http.Handler) http.Handler) *Server {
	log.Printf("add middleware %#v", runtime.FuncForPC(reflect.ValueOf(middleware).Pointer()).Name())
	s.middlewares = append(s.middlewares, middleware)
	return s
}

// AddRoutes adds Routes to the server
//
// param: <routes> is list of Route wrapped in Routes
func (s *Server) AddRoutes(route Route) *Server {
	RouteTable = append(RouteTable, route)
	return s
}

// Build add all the provided configurations to the http.Server
// definition from NewServer or DefaultServer
func (s *Server) Build() *Server {
	s.Handler = addRoutes(RouteTable, s.StrictSlash)

	for _, m := range s.middlewares {
		s.addMiddleware(m)
	}
	return s
}

// Start the http.Server daemon
//
// returns chan os.Signal
func (s *Server) Start() chan os.Signal {
	log.Println("starting server daemon")

	s.Build()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatalf("%s", err.Error())
		}
	}()

	return shutdown
}

// Stop the http.Server daemon
//
// returns shutdown error
func (s *Server) Stop() error {
	log.Println("stopping server daemon")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	return s.Shutdown(ctx)
}

// AddCORSHandler is pre-defined CORS configuration that
// is used to configure CORS without needing for a custom http.Handler
//
// param: <handlerConfig> is HttpResponseConfig definition for CORS config
func (s *Server) AddCORSHandler(handlerConfig HttpResponseConfig) *Server {
	s.AddMiddleware(func(next http.Handler) http.Handler {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", handlerConfig.ContextType)
			w.Header().Set("Access-Control-Allow-Origin", handlerConfig.AccessControlAllowOrigin)
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(handlerConfig.AccessControlAllowMethods, ","))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(handlerConfig.AccessControlAllowHeaders, ","))
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
		return h
	})
	return s
}

func (s *Server) addMiddleware(middleware func(http.Handler) http.Handler) *Server {
	s.Handler = middleware(s.Handler)
	return s
}
