package gre

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

/**
 * Package name: gre
 * Project name: go-routing-engine
 * Created by: Praveen Premaratne
 * Created on: 01/05/2023 14:39
 */

func ExampleDefaultServer() {

	server := DefaultServer(9999, false).
		AddCORSHandler(HttpResponseConfig{
			ContextType:               "application/json",
			AccessControlAllowOrigin:  "*",
			AccessControlAllowMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
			AccessControlAllowHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		}).
		AddRoutes(Route{Name: "Hello",
			Methods:    []string{http.MethodGet},
			Pattern:    "/hello",
			Deprecated: false,
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, "{\"message\": \"hello\"}")
			},
		}).
		Build()

	<-server.Start()

	if err := server.Stop(); err != nil {
		log.Printf("%s", err.Error())
	}

	// Output:
	// 2023/04/29 21:23:45 add middleware "github.com/razorcorp/go-routing-engine/gre.(*Server).AddCORSHandler.func1"
	// 2023/04/29 21:23:45 add global handler 404 - not found
	// 2023/04/29 21:23:45 add global handler 405 - method not allowed
	// 2023/04/29 21:23:45 add mapping: Hello ( [GET] /hello )
	// 2023/04/29 21:23:45 starting server daemon
	// 2023/04/29 21:23:55 stopping server daemon
	// 2023/04/29 21:23:55 http: Server closed

}

func ExampleNewServer() {
	server := &Server{}
	server.Addr = "0.0.0.0:8080"
	server.ReadTimeout = time.Duration(15) * time.Second
	server.WriteTimeout = time.Duration(15) * time.Second
	server.StrictSlash = false
	server.Handler = nil

	//Output:
}

func appRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				_, _ = fmt.Fprintln(os.Stdout, "fetal error, application unexpectedly exited")
				_, _ = fmt.Fprintf(os.Stdout, "%#v", err)
				w.WriteHeader(http.StatusInternalServerError)
				response := ErrorResponse{
					Code:  http.StatusInternalServerError,
					Cause: "oops, something went wrong. we're looking into it",
				}
				fmt.Fprint(w, response.Json())
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func ExampleServer_AddMiddleware() {
	server := DefaultServer(8080, false)
	server.AddRoutes(Route{
		Name:       "Crash",
		Methods:    []string{"GET"},
		Pattern:    "/crash",
		Deprecated: false,
		HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
			panic("fake application crash")
		},
	})

	server.AddMiddleware(appRecovery)

	server.Build()

	<-server.Start()

	if err := server.Stop(); err != nil {
		log.Printf("%#v", err)
	}

	// Output:
	// 2023/05/01 15:11:29 add middleware "main.appRecovery"
	// 2023/05/01 15:11:29 add global handler 404 - not found
	// 2023/05/01 15:11:29 add global handler 405 - method not allowed
	// 2023/05/01 15:11:29 add mapping: Hello ( [GET] /hello )
	// 2023/05/01 15:11:29 add mapping: Crash ( [GET] /crash )
	// 2023/05/01 15:11:29 starting server daemon
	// 2023/05/01 15:11:30  GET /hello Hello 0s PostmanRuntime/7.32.2
	// 2023/05/01 15:11:33  GET /crash Crash 0s PostmanRuntime/7.32.2
	// fetal error, application unexpectedly exited: fake application crash
	// crash recovered
}

func ExampleServer_AddResponseMiddleware() {
	server := DefaultServer(8080, false)
	server.AddMiddleware(appRecovery)
	server.AddRoutes(Route{Name: "Hello",
		Methods:    []string{http.MethodGet},
		Pattern:    "/hello",
		Deprecated: false,
		HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "{\"message\": \"hello\"}")
		},
	})

	server.AddCORSHandler(HttpResponseConfig{
		ContextType:               "",
		AccessControlAllowOrigin:  "",
		AccessControlAllowMethods: nil,
		AccessControlAllowHeaders: nil,
	})

	server.Build()

	<-server.Start()

	if err := server.Stop(); err != nil {
		log.Printf("%#v", err)
	}

	// Output:
	// 2023/05/01 19:37:48 add middleware "github.com/razorcorp/go-routing-engine/gre.appRecovery"
	// 2023/05/01 19:37:48 add middleware "github.com/razorcorp/go-routing-engine/gre.(*Server).AddCORSHandler.func1"
	// 2023/05/01 19:37:48 add global handler 404 - not found
	// 2023/05/01 19:37:48 add global handler 405 - method not allowed
	// 2023/05/01 19:37:48 add mapping: Hello ( [GET] /hello )
}

func ExampleServer_AddRoutes() {

	server := DefaultServer(9999, false)

	server.AddRoutes(Route{Name: "Hello",
		Methods:    []string{http.MethodGet},
		Pattern:    "/hello",
		Deprecated: false,
		HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "{\"message\": \"hello\"}")
		},
	})

	server.Build()

	<-server.Start()

	if err := server.Stop(); err != nil {
		log.Printf("%s", err.Error())
	}

	// Output:
	// 2023/04/29 21:23:45 add middleware "github.com/razorcorp/go-routing-engine/gre.(*Server).AddCORSHandler.func1"
	// 2023/04/29 21:23:45 add global handler 404 - not found
	// 2023/04/29 21:23:45 add global handler 405 - method not allowed
	// 2023/04/29 21:23:45 add mapping: Hello ( [GET] /hello )

}
