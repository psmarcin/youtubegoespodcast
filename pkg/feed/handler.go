package feed

import (
	"net/http"
	"time"
	"ygp/pkg/cache"
	"ygp/pkg/errx"
	"ygp/pkg/utils"

	"github.com/gorilla/mux"
)

const (
	cacheFeedPrefix = "feed_"
	cacheFeedTTL    = time.Hour * 24 * 1
)

func Handler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	channelID := params["channelId"]
	searchPhrase := r.FormValue("search")
	cacheKey := cacheFeedPrefix + channelID + "_" + searchPhrase

	w.Header().Set("Content-Type", "application/rss+xml; charset=UTF-8")

	// get cache
	c, _ := cache.Client.GetKey(cacheKey, nil)

	if c != "" {
		utils.Send(w, c, http.StatusOK)
		return
	}

	f := new(channelID)
	err := f.getDetails(channelID)

	if err.IsError() {
		utils.SendError(w, err)
		return
	}
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

	// set cache
	go cache.Client.SetKey(cacheKey, string(serialized), cacheFeedTTL)

	utils.Send(w, string(serialized), http.StatusOK)
}

func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}
