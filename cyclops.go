package cyclops

import (
	"log"
	"net/http"
)

func StartServer(server *http.Server) {
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

func StartTLSServer(server *http.Server, certFile, keyFile string) {
	if err := server.ListenAndServeTLS(certFile, keyFile); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTPS server ListenAndServe: %v", err)
	}
}