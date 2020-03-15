package api

import (
	"net/http"

	"ygp/pkg/youtube"
)

var disableCache = false

// Handler is a entrypoint for router
func TrendingHandler(w http.ResponseWriter, r *http.Request) {
	response, err := youtube.GetTrending(disableCache)
	if err.IsError() {
		SendError(w, err)
		return
	}
	serialized := youtube.Serialize(response)
	Send(w, serialized, http.StatusOK)
}