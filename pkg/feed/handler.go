package feed

import (
	"net/http"
	"ytg/pkg/db"

	"github.com/gorilla/mux"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	channelID := params["channelId"]

	f := new(channelID)
	err := f.getDetails(channelID)

	// save log to database
	go db.DB.SaveChannel(ctx, channelID, err)

	if err != nil {
		handleError(w, err)
		return
	}
	searchPhrase := r.FormValue("search")
	videos, err := f.getVideos(searchPhrase)
	if err != nil {
		handleError(w, err)
		return
	}
	err = f.setVideos(videos)
	if err != nil {
		handleError(w, err)
		return
	}
	serialized, err := f.serialize()
	if err != nil {
		handleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/rss+xml; charset=UTF-8")
	w.Write([]byte(serialized))
}

func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}
