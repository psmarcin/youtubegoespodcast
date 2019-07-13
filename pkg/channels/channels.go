package channels

import (
	"net/http"
	"ygp/pkg/utils"
	"ygp/pkg/youtube"
)

// Handler is default router handler for GET /channel endpoint
func Handler(w http.ResponseWriter, r *http.Request) {
	q := r.FormValue("q")
	channels, err := youtube.GetChannels(q)
	if err.IsError() {
		utils.SendError(w, err)
		return
	}
	serialied := youtube.Serialize(channels)
	utils.Send(w, serialied, http.StatusOK)
}
