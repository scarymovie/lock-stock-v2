package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingAllRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("Request: Method=%s, URL=%s, RemoteAddr=%s", r.Method, r.URL.Path, r.RemoteAddr)

		ww := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(ww, r)

		duration := time.Since(start)

		log.Printf("Response: Status=%d, Duration=%s", ww.statusCode, duration)
	})
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func NotFoundHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Обработка запроса
		next.ServeHTTP(w, r)

		if ww, ok := w.(*responseWriterWrapper); ok && ww.statusCode == http.StatusOK {
			http.NotFound(w, r)
			log.Printf("404 Not Found: URL=%s", r.URL.Path)
		}
	})
}
