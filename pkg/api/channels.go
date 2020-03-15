package api

import (
	"net/http"

	"ygp/pkg/youtube"
)

// Handler is default router handler for GET /channel endpoint
func ChannelsHandler(w http.ResponseWriter, r *http.Request) {
	q := r.FormValue("q")
	channels, err := youtube.GetChannels(q)
	if err.IsError() {
		SendError(w, err)
		return
	}
	serialied := youtube.Serialize(channels)
	Send(w, serialied, http.StatusOK)
}
