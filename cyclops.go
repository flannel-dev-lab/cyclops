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

	http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
}
