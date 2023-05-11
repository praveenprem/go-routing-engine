package gre

import (
	"log"
)

/**
 * Package name: gre
 * Project name: go-routing-engine
 * Created by: Praveen Premaratne
 * Created on: 01/05/2023 12:49
 */

func ExampleNewRouter() {
	routes := Routes{
		Route{
			Name:        "Hello",
			Methods:     []string{"GET", "POST"},
			Pattern:     "/hello",
			Deprecated:  false,
			HandlerFunc: nil,
		},
	}
	router := NewRouter(routes, false)

	log.Printf("%#v", router)

	// Output:
	// 2023/05/01 13:02:45 add global handler 404 - not found
	// 2023/05/01 13:02:45 add global handler 405 - method not allowed
	// 2023/05/01 13:02:45 add mapping: Hello ( [GET POST] /hello )
	// 2023/05/01 13:02:45 &mux.Router{NotFoundHandler:(http.HandlerFunc)(0xc97580), MethodNotAllowedHandler:(http.HandlerFunc)(0xc97680),
	// routes:[]*mux.Route{(*mux.Route)(0xc0002141e0), (*mux.Route)(0xc0002143c0)}, namedRoutes:map[string]*mux.Route{"Hello":(*mux.Route)(0xc0002141e0)},
	// KeepContext:false, middlewares:[]mux.middleware{(mux.MiddlewareFunc)(0xc76c80), (mux.MiddlewareFunc)(0xc96b80)},
	// routeConf:mux.routeConf{useEncodedPath:false, strictSlash:false, skipClean:false, regexp:mux.routeRegexpGroup{host:(*mux.routeRegexp)(nil),
	// path:(*mux.routeRegexp)(nil), queries:[]*mux.routeRegexp(nil)}, matchers:[]mux.matcher(nil), buildScheme:"", buildVarsFunc:(mux.BuildVarsFunc)(nil)}}
}
