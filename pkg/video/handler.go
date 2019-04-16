package video

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"ytg/pkg/utils"
)

// Handler handles GET /video?videoId= route
func Handler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("[API] request %s %s", r.Method, r.RequestURI)
	videoID := r.FormValue("videoId")
	videoURL := GetURL(videoID)
	if videoURL == "" {
		http.NotFound(w, r)
		return
	}
	utils.Redirect(w, videoURL)
}
