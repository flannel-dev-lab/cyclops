// cyclops is a minimal web framework
package cyclops

import (
	"fmt"
	"github.com/flannel-dev-lab/cyclops/input"
	"log"
	"net/http"
	"strings"
)

const banner = `
 ______ ______/ /__  ___  ___
/ __/ // / __/ / _ \/ _ \(_-<
\__/\_, /\__/_/\___/ .__/___/
   /___/          /_/

https://github.com/flannel-dev-lab/cyclops
`

// StartServer starts a simple http server
func StartServer(address string, handler http.Handler) {
	fmt.Print(banner)
	if err := http.ListenAndServe(address, handler); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

// StartTLSServer starts a TLS server with provided TLS cert and key files
func StartTLSServer(address string, handler http.Handler, certFile, keyFile string) {
	fmt.Print(banner)

	if err := http.ListenAndServeTLS(address, certFile, keyFile, handler); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTPS server ListenAndServe: %v", err)
	}
}

// Param - Get a url parameter by name
func Param(r *http.Request, name string) string {
	return input.Query(":"+name, r)
}

// ParamNames - Get a url parameter name list with the leading :
func ParamNames(r *http.Request) []string {
	var names []string
	for k := range r.URL.Query() {
		if strings.HasPrefix(k, ":") {
			names = append(names, k)
		}
	}
	return names
}
