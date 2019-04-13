package utils

import (
	"io"
	"net/http"
)

// JSONResponse set headers for json response
func JSONResponse(w http.ResponseWriter) {
	w.Header().Set("content-type", "application/json")
}

// OkResponse set statusCode for response
func OkResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

// WriteBodyResponse writes body to writer
func WriteBodyResponse(w http.ResponseWriter, body string) {
	io.WriteString(w, body)
}

func AllowCorsResponse(w http.ResponseWriter, r *http.Request) {
	switch origin := r.Header.Get("Origin"); origin {
	case "http://localhost:8080", "https://yt.psmarcin.dev", "https://yt.psmarcin.me":
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
}
