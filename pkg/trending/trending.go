package trending

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"ytg/pkg/utils"
	"ytg/pkg/youtube"
)

// Handler is a entrypoint for router
func Handler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("[API] request %s %s", r.Method, r.RequestURI)
	response := youtube.GetTrendings()
	serialied := youtube.SerializeTrending(response)
	utils.JSONResponse(w)
	utils.AllowCorsResponse(w, r)
	utils.OkResponse(w)
	utils.WriteBodyResponse(w, serialied)
}
