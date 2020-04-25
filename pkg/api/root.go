package api

import "net/http"

func rootHandler(w http.ResponseWriter, r *http.Request) {
	JSONResponse(w)
	Send(w, RootJSON, http.StatusOK)
}
