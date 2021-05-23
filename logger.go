package main

import (
	"log"
	"net/http"
)

type loggerResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *loggerResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func middlewareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lw := &loggerResponseWriter{w, http.StatusOK}
		next.ServeHTTP(lw, r)
		log.Printf("%s %s: %d", r.Method, r.RequestURI, lw.statusCode)
	})
}
