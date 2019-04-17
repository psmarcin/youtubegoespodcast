package trending

import (
	"net/http"
	"ytg/pkg/utils"
	"ytg/pkg/youtube"
)

// Handler is a entrypoint for router
func Handler(w http.ResponseWriter, r *http.Request) {
	response := youtube.GetTrendings()
	serialied := youtube.Serialize(response)
	utils.OkResponse(w)
	utils.WriteBodyResponse(w, serialied)
}
