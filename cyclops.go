// cyclops is a minimal web framework
package cyclops

import (
	"log"
	"net/http"
)

// StartServer starts a simple http server
func StartServer(server *http.Server) {
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

// StartTLSServer starts a TLS server with provided TLS cert and key files
func StartTLSServer(server *http.Server, certFile, keyFile string) {
	if err := server.ListenAndServeTLS(certFile, keyFile); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTPS server ListenAndServe: %v", err)
	}
}