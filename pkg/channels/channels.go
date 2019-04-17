package channels

import (
	"net/http"
	"ytg/pkg/utils"
	"ytg/pkg/youtube"
)

// Handler is default router handler for GET /channel endpoint
func Handler(w http.ResponseWriter, r *http.Request) {
	q := r.FormValue("q")
	channels := youtube.GetChannels(q)
	serialied := youtube.Serialize(channels)
	utils.OkResponse(w)
	utils.WriteBodyResponse(w, serialied)
}
