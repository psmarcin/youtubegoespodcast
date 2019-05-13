package utils

import "net/http"

// MiddlewareJSON set header to response with JSON
func MiddlewareJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		JSONResponse(w)
		next.ServeHTTP(w, r)
	})
}

// MiddlewareCORS set header to proper CORS headers
func MiddlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		AllowCorsResponse(w, r)
		next.ServeHTTP(w, r)
	})
}
