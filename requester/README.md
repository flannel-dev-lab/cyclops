# Requestor
- Requestor is a feature of cyclops which make handling http requests easy as a pie.

- A HTTP GET using cyclops requester
```go
package main

import (
	"github.com/flannel-dev-lab/cyclops/requester"
	"net/http"
	"log"
	"io/ioutil"
)

func main() {
	resp, err := requester.Get("https://httpbin.org/get", nil, nil)
	if err != nil {
		panic(err)
	}
	
	defer resp.Body.Close()
}
```
