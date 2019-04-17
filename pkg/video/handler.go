package video

import (
	"github.com/gorilla/mux"
	"net/http"
	"ytg/pkg/utils"
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
	utils.Redirect(w, videoURL)
}

// RedirectHandler is legacy handler for old routes, it redirect to new route
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	utils.PermanentRedirect(w, "/video/"+r.FormValue("videoId")+".mp3")
}
