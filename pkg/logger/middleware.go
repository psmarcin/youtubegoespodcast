package logger

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// Middleware act as middleware and log request method and path
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("[API] Request %s %s %+v", r.Method, r.RequestURI, r.Header)
		next.ServeHTTP(w, r)
	})
}
