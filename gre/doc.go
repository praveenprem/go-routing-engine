/*
Package gre implements both router/dispatcher and web server.

The gre package within provide both a mux.Router as we all a http.Server
for users looking for advanced capabilities of a strip down HTTP request multiplexer
and for those who looking for a basic HTTP server that have predefined methods for
adding rotes and middleware. The main features are:

  - Set gre.Route structure for registering routes
  - Router and dispatcher built using github.com/gorilla/mux
  - Preconfigured http.Server option for hassle-free deployment
  - Prometheus metrics built in with auto register to allow customer metrics registration

Let's start building a simple HTTP server:

	// Following implementation of the routes allow users to scatter route definition across multiple files in cases such as large api implementations

	func main() {

		gre.RouteTable = append(gre.RouteTable, gre.Route{Name: "Hello",
			Methods:     []string{http.MethodGet},
			Pattern:     "/hello",
			Deprecated:  false,
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, "{\"message\":\"Hello\"}")
			},
		})

		server := gre.DefaultServer(9999, false).
			AddCORSHandler(gre.HttpResponseConfig{
				ContextType:               "application/json",
				AccessControlAllowOrigin:  "*",
				AccessControlAllowMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
				AccessControlAllowHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
			}).
			Build()

		<-server.Start()

		if err := server.Stop(); err != nil {
			log.Printf("%s", err.Error())
		}
	}

Here we register a gre.DefaultServer with a middleware and a route. Server will serve on port 9999
and strict slashes to false (if registered rote is `/path`, accessing `/path/` will not match).

	server := gre.DefaultServer(9999, false)

Then we add a predefined middleware that sets the response structure for CORS queries. Following
that we've the route definition and the request handler function and its code. This can be a reference
to an actual function as long as the function resembles the http.HandlerFunc pattern.

	server.
		AddCORSHandler(gre.HttpResponseConfig{
			ContextType:               "application/json",
			AccessControlAllowOrigin:  "*",
			AccessControlAllowMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
			AccessControlAllowHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		}).
		AddRoutes(gre.Routes{
			gre.Route{Name:       "Hello",
				Methods:    []string{http.MethodGet},
				Pattern:    "/hello",
				Deprecated: false,
				HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					fmt.Fprint(w, "{\"message\": \"hello\"}")
				},
			},
		})

Build() method must be invoked when using the gre.Server as this will add the routes and the middleware
to the mux.Router and starts the http.Server and returns a channel to monitoring for os.Interrupt,
syscall.SIGINT, syscall.SIGTERM which is used to invoke Stop() method to terminate the server gracefully.

	server.Build()

	<-server.Start()

	if err := server.Stop(); err != nil {
		log.Printf("%s", err.Error())
	}
*/
package gre
