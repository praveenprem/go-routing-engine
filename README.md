# Go Routing Engine

[![stable](https://img.shields.io/github/v/tag/razorcorp/go-routing-engine?sort=semver&logo=github&style=for-the-badge&label=Tag)](https://github.com/razorcorp/go-routing-engine/releases/latest)
![GitHub](https://img.shields.io/github/license/razorcorp/go-routing-engine?style=for-the-badge)

[//]: # (![GitHub Release Date]&#40;https://img.shields.io/github/tag-date/razorcorp/go-routing-engine?display_date=published_at&style=for-the-badge&label=Release%20Date&#41;)

GRE is a light framework for building RestAPIs for small and large scale applications. This is a simple framework build on 
router and dispatcher.

Unlike other popular frameworks, this framework is less opinionated and give user more control over the implementation.
However, this framework does have few restrictions:
- It implements a specific structure for registering routes and handlers
- Has a predefined JSON error response structure for common errors. I.E: 404 and 405
- Has a set API request logger
- Has Prometheus metrics builtin

> More information on route options, refer to https://pkg.go.dev/github.com/gorilla/mux#Route.Path

## Installation
[![Go Reference](https://pkg.go.dev/badge/github.com/razorcorp/go-routing-engine/gre.svg)](https://pkg.go.dev/github.com/razorcorp/go-routing-engine/gre)
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/razorcorp/go-routing-engine/master?filename=go.mod&&style=flat)

With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain:

```shell
go get github.com/razorcorp/go-routing-engine
```

## Example

Let's start building a simple HTTP server:

> Following implementation of the routes allow users to scatter route definition across multiple files in cases such as
> large api implementations

```go
package main

import (
	"fmt"
	"github.com/razorcorp/go-routing-engine/gre"
	"log"
	"net/http"
)

func main() {

	gre.RouteTable = append(gre.RouteTable, gre.Route{Name: "Hello",
		Methods:    []string{http.MethodGet},
		Pattern:    "/hello",
		Deprecated: false,
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
```

Here we register a new route on the package RouteTable array, followed by initialising a basic HTTP server. This in an
instance of `http.Server` and be further customised with your own configurations or use other pre-defined methods to add
additional functionality. Such as the middleware implementation of CORS in the example above.

#### Note
> When using builtin methods, `.Build()` method must be called at the end of the configuration chain, as all method 
configurations won't take effect until this method is called.

---

More examples can be found in the [example](example) directory.
