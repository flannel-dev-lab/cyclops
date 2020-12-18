# Cyclops
[![Go Report Card](https://goreportcard.com/badge/github.com/flannel-dev-lab/cyclops)](https://goreportcard.com/report/github.com/flannel-dev-lab/cyclops)
[![Build Status](https://travis-ci.org/flannel-dev-lab/cyclops.svg?branch=master)](https://travis-ci.org/flannel-dev-lab/cyclops)
[![Go Doc](https://godoc.org/github.com/flannel-dev-lab/cyclops?status.svg)](https://godoc.org/github.com/flannel-dev-lab/cyclops)
[![codecov](https://codecov.io/gh/flannel-dev-lab/cyclops/branch/master/graph/badge.svg)](https://codecov.io/gh/flannel-dev-lab/cyclops)

![GitHub Logo1](Gopher-logo.png)

## Download
```shell
go get github.com/flannel-dev-lab/cyclops/v2
```

## Features
- Plug and Play Middleware support
- Customized response messages

## Table of contents
1. [Middleware](middleware/README.md)
2. [Cookies](cookie/README.md)
3. [Handling Inputs](input/README.md)
4. [Handling HTTP Requests](requester/README.md)
5. [Handling Sessions](sessions/README.md)
6. [Logging](logger/README.md)

## Usage

```go
package main

import (
	"fmt"
	"github.com/flannel-dev-lab/cyclops/v2"
	"github.com/flannel-dev-lab/cyclops/v2/middleware"
	"github.com/flannel-dev-lab/cyclops/v2/response"
	"github.com/flannel-dev-lab/cyclops/v2/router"
	"net/http"
)

func main() {
	routerObj := router.New(false, nil, nil)

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
```
