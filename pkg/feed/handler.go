package feed

import (
	"net/http"
	"ygp/pkg/db"
	"ygp/pkg/errx"
	"ygp/pkg/utils"

	"github.com/gorilla/mux"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	channelID := params["channelId"]

	f := new(channelID)
	err := f.getDetails(channelID)

	// save log to database
	go db.DB.SaveChannel(ctx, channelID, err.Err)

	if err.IsError() {
		utils.SendError(w, err)
		return
	}
	searchPhrase := r.FormValue("search")
	videos, getVideoErr := f.getVideos(searchPhrase)
	if getVideoErr.IsError() {
		utils.SendError(w, getVideoErr)
		return
	}

	setVideosErr := f.setVideos(videos)
	if setVideosErr.IsError() {
		utils.SendError(w, setVideosErr)
		return
	}
	serialized, serializeErr := f.serialize()
	if serializeErr != nil {
		utils.SendError(w, errx.NewAPIError(serializeErr, http.StatusInternalServerError))
		return
	}
	w.Header().Set("Content-Type", "application/rss+xml; charset=UTF-8")
	utils.Send(w, string(serialized), http.StatusOK)
}

func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}
