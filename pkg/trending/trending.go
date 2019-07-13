package trending

import (
	"net/http"
	"ygp/pkg/utils"
	"ygp/pkg/youtube"
)

// Handler is a entrypoint for router
func Handler(w http.ResponseWriter, r *http.Request) {
	response, err := youtube.GetTrending()
	if err.IsError() {
		utils.SendError(w, err)
		return
	}
	serialied := youtube.Serialize(response)
	utils.Send(w, serialied, http.StatusOK)
}
