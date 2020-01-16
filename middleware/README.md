# Middlewares

- Cyclops supports middlewares and few middlewares come pre-included with the project and they are:
    - CORS
    - Panic Handler
    - Request Logger
    - Set Secure Headers with every request
    - Set Default Headers with every request

- Out of the above middlewares, PanicHandler and RequestLogger are enabled by default
- You can also add your own custom middlewares, the only  thing to take care of when writing custom middlewares is that
the function should take in `http.Handler` as a parameter and return `http.Handler`

### Using CORS Middleware
Cyclops supports CORS and can be used as explained below

```
package main

import (
	"fmt"
	"github.com/flannel-dev-lab/cyclops"
	"github.com/flannel-dev-lab/cyclops/middleware"
	"github.com/flannel-dev-lab/cyclops/router"
	"html/template"
	"net/http"
	"time"
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

	routerObj := router.New()
	routerObj.Get("/", middleware.NewChain(cors.CORSHandler).Then(Login))

	cyclops.StartServer(":8080", routerObj)
}

func Login(w http.ResponseWriter, r *http.Request) {
	response.SuccessResponse(200, w, nil)
}

```

### Setting Default Headers with every request
```
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

	~~defaultHeaders := middleware.DefaultHeaders{
        ContentType: "application/json",
	}~~

    routerObj := router.New()
    routerObj.Get("/", middleware.NewChain(DefaultHeaders.DefaultHeaders).Then(Login))
    
    cyclops.StartServer(":8080", routerObj)
}

func Login(w http.ResponseWriter, r *http.Request) {
	response.SuccessResponse(200, w, nil)
}

```
When we run the above code, the response header `Content-Type: application/json` is set with every request

## Middleware Chaining
If you want to use multiple middlewares for a request, cyclops allows you to do that as well. All you need to do is like
below:
```
middleware.NewChain(defaultHeaders.DefaultHeaders, cors.CORSHandler).Then(http.HandlerFunc(Login))
```