package video

import (
	"net/http"
	"ygp/pkg/utils"

	"github.com/gorilla/mux"
)

// Handler handles GET /video?videoId= route
func Handler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	videoID := params["videoId"]

	videoURL := GetURL(videoID)
	if videoURL == "" {
		http.NotFound(w, r)
		return
	}
	w.Header().Del("Content-Type")
	utils.Redirect(w, videoURL)
}

// RedirectHandler is legacy handler for old routes, it redirect to new route
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	utils.PermanentRedirect(w, "/video/"+r.FormValue("videoId")+".mp3")
}
