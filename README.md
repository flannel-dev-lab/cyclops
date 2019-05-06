# Cyclops

## Alerting
Alerting is built into cyclops and currently supports [Sentry](https://sentry.io).
Alerting guide can be found [here](alerts/README.md)

## Features
- Plug and Play Middleware support

## Middlewares

### Panic Handler
- Handles server crashes gracefully, Sends a 500 return code. Enabled automatically

### Request Logger
- Logs access requests. Enabled automatically
- `2019/05/05 19:58:40 [::1]:58225 POST /hello localhost:8080 HTTP/1.1`

### Set Response Headers
- Sets default response headers. The following response headers are set  by default
```
w.Header().Set("X-Content-Type-Options", "nosniff")
w.Header().Set("X-Frame-Options", "deny")
w.Header().Set("Content-Type", "application/json")
w.Header().Set("X-XSS-Protection", "1; mode=block")
```
- The above headers can be overridden by passing the headers in your own handlers
- Setting middleware is as easy as pie
```
// Add middlewares to a route that you define
routes["/hello"] = middleware.NewChain(middleware.SetHeaders).Then(http.HandlerFunc(Root))
```
- The method `NewChain()` will accept multiple middlewares  of the form `func(http.Handler) http.Handler`

## Usage
```
package main

import (
	"WebFramework/middleware"
	"WebFramework/router"
	"fmt"
	"log"
	"net/http"
)

func Root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Root!\n")
}


func main() {
	handler, server := router.InitializeServer(":8080")

	routes := make(map[string]http.Handler)

	routes["/hello"] = middleware.NewChain(middleware.RequestLogger, middleware.SetHeaders).Then(http.HandlerFunc(Root))
	routes["/bye"] = middleware.NewChain(middleware.RequestLogger).Then(http.HandlerFunc(Root))



	router.RegisterRoutes(handler, routes)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

```

## Testing
`go test ./...`

** FEATURES **
- Middleware support
- Any case route paths
- Panic handle Middlewares