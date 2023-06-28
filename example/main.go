package main

import (
	"encoding/json"
	"fmt"
	"github.com/razorcorp/go-routing-engine/gre"
	"log"
	"net/http"
	"os"
)

/**
 * Package name: main
 * Project name: go-routing-engine
 * Created by: Praveen Premaratne
 * Created on: 27/04/2023 23:11
 */

type Response struct {
	Message string `json:"message"`
}

//var totalRequests = prometheus.NewCounterVec(
//	prometheus.CounterOpts{
//		Name: "http_requests_total",
//		Help: "Number of get requests.",
//	},
//	[]string{"path"},
//)

func appRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "fetal error, application unexpectedly exited: %s\n", err)
				_, _ = fmt.Println("crash recovered")
				w.WriteHeader(http.StatusInternalServerError)
				response := gre.ErrorResponse{
					Code:  http.StatusInternalServerError,
					Cause: "oops, something went wrong. we're looking into it",
				}
				fmt.Fprint(w, response.Json())
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {

	gre.RouteTable = append(gre.RouteTable, gre.Route{Name: "Hello",
		Methods:     []string{http.MethodGet},
		Pattern:     "/hello",
		Deprecated:  false,
		HandlerFunc: hello,
	})

	server := gre.DefaultServer(9999, false).
		AddCORSHandler(gre.HttpResponseConfig{
			ContextType:               "application/json",
			AccessControlAllowOrigin:  "*",
			AccessControlAllowMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
			AccessControlAllowHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		}).
		AddMiddleware(appRecovery)
	//AddRoutes(gre.Route{Name: "Hello",
	//	Methods:     []string{http.MethodGet},
	//	Pattern:     "/hello",
	//	Deprecated:  false,
	//	HandlerFunc: hello,
	//}).

	<-server.Start()

	if err := server.Stop(); err != nil {
		log.Printf("%s", err.Error())
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	resp := &Response{Message: "Hello"}

	//time.Sleep(time.Duration(rand.Intn(10)) * time.Second)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, resp.Json())
}

func (r *Response) Json() string {
	body, err := json.Marshal(r)
	if err != nil {
		log.Printf("%#v", r)
		response := gre.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Cause: "something went wrong, try again in few minutes",
			Debug: "err: JSON encoding failed",
		}
		return response.Json()
	}
	return string(body)
}
