package gre

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

/**
 * Package name: gre
 * Project name: go-routing-engine
 * Created by: Praveen Premaratne
 * Created on: 27/04/2023 21:25
 */

type (
	// Routes defines a collection of Route
	Routes []Route

	// Route defines parameters used to configure mux.Router
	Route struct {

		// Name of the route used for request logging
		Name string

		// Methods are list of HTTP methods allowed for this route
		Methods []string

		// Pattern is the URI path pattern
		// This can include variables if path contains dynamic segments
		// example: "/user/{name}"
		Pattern string

		// Deprecated allows a route to be flagged without
		// completely removing it from code any route set
		// as deprecated will respond with defined error response
		Deprecated bool

		// HandlerFunc is an adapter to allow the use of
		// ordinary functions as HTTP handlers.
		HandlerFunc http.HandlerFunc
	}

	// Server extends http.Server with few additional parameters
	Server struct {

		// StrictSlash defines the trailing slash behavior for new routes
		StrictSlash bool

		// middlewares is an internal component for adding middleware to
		// router globally
		middlewares []func(http.Handler) http.Handler

		http.Server
	}

	// HttpResponseConfig is built in CORS configurator
	HttpResponseConfig struct {

		// ContextType is the HTTP response format
		ContextType string

		// AccessControlAllowOrigin is CORS preflight check for restricting
		// API request source
		AccessControlAllowOrigin string

		// AccessControlAllowMethods is CORS preflight check for HTTP methods
		// api requests allowed to use
		AccessControlAllowMethods []string

		// AccessControlAllowHeaders CORS preflight check for indicating which
		// HTTP headers are allowed
		AccessControlAllowHeaders []string
	}

	// ErrorResponse defines the default error response structure for HTTP requests
	ErrorResponse struct {

		// Code is HTTP status code
		Code int `json:"code"`

		// Cause is the user-friendly error message
		Cause string `json:"cause"`

		// Debug is for additional information if needed
		Debug string `json:"debug,omitempty"`

		looping bool
	}
)

// Json encode ErrorResponse object for transmit.
// returns a string
func (e *ErrorResponse) Json() string {
	body, err := json.Marshal(e)
	if err != nil {
		if e.looping {
			log.Fatalln("error response looping detected")
		}
		log.Printf("%#v", e)
		response := ErrorResponse{
			Code:    http.StatusInternalServerError,
			Cause:   "something went wrong, try again in few minutes",
			Debug:   "err: JSON encoding failed",
			looping: true,
		}
		return response.Json()
	}
	return string(body)
}

func (r *Routes) conflict() error {
	if RouteTable != nil {
		return errors.New("AddRoute method not allowed while RouteTable is set")
	}

	return nil
}
