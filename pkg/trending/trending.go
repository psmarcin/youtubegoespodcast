package trending

import (
	"net/http"
	"ytg/pkg/utils"
	"ytg/pkg/youtube"
)

// Handler is a entrypoint for router
func Handler(w http.ResponseWriter, r *http.Request) {
	response, err := youtube.GetTrendings()
	if err.IsError() {
		utils.SendError(w, err)
		return
	}
	serialied := youtube.Serialize(response)
	utils.Send(w, serialied, http.StatusOK)
}
