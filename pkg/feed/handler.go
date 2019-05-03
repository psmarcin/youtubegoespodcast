package feed

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	channelID := params["channelId"]
	f := new(channelID)
	err := f.getDetails(channelID)
	if err != nil {
		handleError(w, err)
		return
	}
	videos, err := f.getVideos()
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
