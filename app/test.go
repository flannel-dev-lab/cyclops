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
	routerObj := router.New()

	routerObj.Get("/hello", middleware.NewChain().Then(Hello))
	routerObj.Post("/bye", Bye)

	// Named parameters can be used as
	routerObj.Get("/users/:name", middleware.NewChain().Then(PathParam))

	// static can be registered as
	routerObj.RegisterStatic("{PATH TO STATIC DIRECTORY}", "/static/")
	cyclops.StartServer(":8080", routerObj)
}

func PathParam(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi %s", cyclops.Param(r, "name"))
	response.SuccessResponse(200, w, nil)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	response.SuccessResponse(200, w, nil)
}

func Bye(w http.ResponseWriter, r *http.Request) {
	response.SuccessResponse(200, w, nil)
}