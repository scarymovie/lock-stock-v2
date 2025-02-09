package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingAllRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Upgrade") == "websocket" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		ww := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		log.Printf("Before handler: URL=%s, Method=%s, Initial Status=%d", r.URL.Path, r.Method, ww.statusCode)

		next.ServeHTTP(ww, r)

		duration := time.Since(start)

		log.Printf("After handler: URL=%s, Method=%s, Final Status=%d, Duration=%s",
			r.URL.Path, r.Method, ww.statusCode, duration)
	})
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(statusCode int) {
	if rw.statusCode != http.StatusOK {
		return
	}
	rw.statusCode = statusCode
	log.Printf("WriteHeader called with status: %d", statusCode)
	rw.ResponseWriter.WriteHeader(statusCode)
}
