package utils

import (
	"net/http"
	"ytg/pkg/errx"

	"github.com/sirupsen/logrus"
)

// JSONResponse set headers for json response
func JSONResponse(w http.ResponseWriter) {
	w.Header().Set("content-type", "application/json")
}

// AllowCorsResponse set proper CORS headers
func AllowCorsResponse(w http.ResponseWriter, r *http.Request) {
	switch origin := r.Header.Get("Origin"); origin {
	case "http://localhost:8080", "https://yt.psmarcin.dev", "https://yt.psmarcin.me":
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
}

// Redirect set headers and statusCode to sent redirect response
func Redirect(w http.ResponseWriter, location string) {
	w.Header().Set("location", location)
	w.WriteHeader(http.StatusFound)
}

// PermanentRedirect set headers and statusCode to sent redirect response
func PermanentRedirect(w http.ResponseWriter, location string) {
	w.Header().Set("location", location)
	w.WriteHeader(http.StatusPermanentRedirect)
}

// Send check if error exists if so return serialized version if it doesn't it return conetent
func Send(w http.ResponseWriter, body string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(body))
}

// Send check if error exists if so return serialized version if it doesn't it return conetent
func SendError(w http.ResponseWriter, err errx.APIError) {
	logrus.WithError(err.Err).Printf("[API] Error response")
	w.WriteHeader(err.StatusCode)
	w.Write([]byte(err.Serialize()))
}
