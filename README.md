<p align="center"><img src="![GitHub Logo1](Gopher-logo.png)" /></p>
<p align="center">
Hello, this is cyclops, I am a minimal web framework
[![Go Report Card](https://goreportcard.com/badge/github.com/flannel-dev-lab/cyclops)](https://goreportcard.com/report/github.com/flannel-dev-lab/cyclops)
[![Build Status](https://travis-ci.org/flannel-dev-lab/cyclops.svg?branch=master)](https://travis-ci.org/flannel-dev-lab/cyclops)
[![Go Doc](https://godoc.org/github.com/flannel-dev-lab/cyclops?status.svg)](https://godoc.org/github.com/flannel-dev-lab/cyclops)
</p>

# Cyclops

## Features
- Plug and Play Middleware support
- Plug and Play Alerting support
- Customized response messages

## Table of contents
1. [Alerting](alerts/README.md)
2. [Middleware](middleware/README.md)

## Usage
```go
package main

import (
	"fmt"
	"github.com/flannel-dev-lab/cyclops"
	"github.com/flannel-dev-lab/cyclops/response"
	"github.com/flannel-dev-lab/cyclops/router"
	"github.com/flannel-dev-lab/cyclops/middleware"
	"net/http"
)

func main() {

	routes := make(map[string]http.Handler)
	routes["/hello"] = middleware.NewChain().Then(http.HandlerFunc(Hello))
	routes["/bye"] = middleware.NewChain().Then(http.HandlerFunc(Bye))

	handler, server := router.InitializeServer(":8080")

	router.RegisterRoutes(handler, routes)

	cyclops.StartServer(server)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	response.SuccessResponse(200, w, nil)
}

func Bye(w http.ResponseWriter, r *http.Request) {
	response.SuccessResponse(200, w, nil)
}


```
