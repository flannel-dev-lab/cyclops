// cyclops is a minimal web framework
package cyclops

import (
	"fmt"
	"log"
	"net/http"
)

const banner  = `
 ______ ______/ /__  ___  ___
/ __/ // / __/ / _ \/ _ \(_-<
\__/\_, /\__/_/\___/ .__/___/
   /___/          /_/

https://github.com/flannel-dev-lab/cyclops
`

// StartServer starts a simple http server
func StartServer(server *http.Server) {
	fmt.Print(banner)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

// StartTLSServer starts a TLS server with provided TLS cert and key files
func StartTLSServer(server *http.Server, certFile, keyFile string) {
	fmt.Print(banner)
	if err := server.ListenAndServeTLS(certFile, keyFile); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTPS server ListenAndServe: %v", err)
	}
}
