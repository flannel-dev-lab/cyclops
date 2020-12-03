package middleware

import (
	"errors"
	"fmt"
	"github.com/flannel-dev-lab/cyclops/v2/logger"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// loggingResponseWriter is a custom implementation of http.ResponseWriter to log status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	message    string
}

// NewLoggingResponseWriter Creates a reference loggingResponseWriter
func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK, ""}
}

// WriteHeader takes in a http status code and adds to loggingResponseWriter
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Header maps to the http.ResponseWriter's Header() method
func (lrw *loggingResponseWriter) Header() http.Header {
	return lrw.ResponseWriter.Header()
}

// Write takes in a byte array and writes to response writer
func (lrw *loggingResponseWriter) Write(data []byte) (int, error) {
	lrw.message = string(data)
	return lrw.ResponseWriter.Write(data)
}

// accessLogger is used to log access logs for discover service
func AccessLogger(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		ctx = logger.AddKey(ctx, "timestamp", time.Now().UTC().Format(time.RFC3339))
		ctx = logger.AddKey(ctx, "remote_address", r.RemoteAddr)
		ctx = logger.AddKey(ctx, "method", r.Method)
		ctx = logger.AddKey(ctx, "protocol", r.Proto)
		ctx = logger.AddKey(ctx, "path", r.URL.Path)

		u, err := uuid.NewUUID()
		if err != nil {
			logger.Error(ctx, "could not generate request-id", err)
		} else {
			ctx = logger.AddKey(ctx, "api-request-id", u.String())
		}

		startTime := time.Now().UTC()

		r = r.WithContext(ctx)

		lrw := NewLoggingResponseWriter(w)

		h.ServeHTTP(lrw, r)

		ctx = logger.AddKey(ctx, "status_code", fmt.Sprintf("%d", lrw.statusCode))
		ctx = logger.AddKey(ctx, "duration", fmt.Sprintf("%d", time.Since(startTime).Milliseconds()))

		if lrw.statusCode >= 200 && lrw.statusCode <= 399 {
			logger.Info(ctx, "access_log")
		} else {
			logger.Error(ctx, lrw.message, errors.New(lrw.message))
		}

	}
}
