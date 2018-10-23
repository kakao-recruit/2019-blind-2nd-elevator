package myauth

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

type AuthMiddleware struct {
}

var (
	loc *time.Location
	EndOfTest time.Time
)

func (a *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/viewer") {
			next.ServeHTTP(w, r)
			return
		}

		var token string
		if strings.HasPrefix(r.URL.Path, "/start") == false {
			token = r.Header.Get("X-Auth-Token")
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		start := time.Now()

		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)

		log.WithFields(log.Fields{
			"token":      token,
			"path":       r.URL.Path,
			"elapsed":    time.Since(start).Seconds() * 1000,
			"statuscode": lrw.statusCode,
		}).Debug()
	})
}
