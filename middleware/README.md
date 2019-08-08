# Middleware

## CORS Middleware
Cyclops supports CORS and can be used as explained below

```go
package main

import (
	"fmt"
	"github.com/flannel-dev-lab/cyclops"
	"github.com/flannel-dev-lab/cyclops/middleware"
	"github.com/flannel-dev-lab/cyclops/response"
	"github.com/flannel-dev-lab/cyclops/router"
	"net/http"
)

func main() {

	cors := middleware.CORS{
		AllowedOrigin: "https://www.admin.yombu.com",
		AllowedHeaders: []string{"Content-Type", "referrer", "referrer-type"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedCredentials: true,
		ExposedHeaders: []string{"GET"},
		MaxAge: 100,
	}

	routes := make(map[string]http.Handler)
	routes["/"] = middleware.NewChain(cors.CORSHandler).Then(http.HandlerFunc(test))

	handler, server := router.InitializeHTTPServer(":8080")

	router.RegisterRoutes(handler, routes)
	cyclops.StartServer(server)
}

func test(w http.ResponseWriter, r *http.Request) {
	response.SuccessResponse(200, w, nil)
}

```