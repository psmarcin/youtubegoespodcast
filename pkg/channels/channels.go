package channels

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"ytg/pkg/utils"
	"ytg/pkg/youtube"
)

// Handler is default router handler for GET /channel endpoint
func Handler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("[API] request %s %s", r.Method, r.RequestURI)
	q := r.FormValue("q")
	channels := youtube.GetChannels(q)
	serialied := youtube.Serialize(channels)
	utils.JSONResponse(w)
	utils.AllowCorsResponse(w, r)
	utils.OkResponse(w)
	utils.WriteBodyResponse(w, serialied)
}
