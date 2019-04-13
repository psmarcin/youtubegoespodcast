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
