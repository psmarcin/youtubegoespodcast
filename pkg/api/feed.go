package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"ygp/pkg/cache"
	"ygp/pkg/errx"
	"ygp/pkg/feed"
)

const (
	cacheFeedPrefix = "feed_"
	cacheFeedTTL    = time.Hour * 24 * 1
)

func FeedHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	channelID := params["channelId"]
	searchPhrase := r.FormValue("search")
	cacheKey := cacheFeedPrefix + channelID + "_" + searchPhrase

	w.Header().Set("Content-Type", "application/rss+xml; charset=UTF-8")

	// get cache
	c, _ := cache.Client.GetKey(cacheKey, nil)

	if c != "" {
		Send(w, c, http.StatusOK)
		return
	}

	f := feed.New(channelID)
	err := f.GetDetails(channelID)

	if err.IsError() {
		SendError(w, err)
		return
	}
	videos, getVideoErr := f.GetVideos(searchPhrase, false)
	if getVideoErr.IsError() {
		SendError(w, getVideoErr)
		return
	}

	setVideosErr := f.SetVideos(videos)
	if setVideosErr.IsError() {
		SendError(w, setVideosErr)
		return
	}
	serialized, serializeErr := f.Serialize()
	if serializeErr != nil {
		SendError(w, errx.New(serializeErr, http.StatusInternalServerError))
		return
	}

	// set cache
	go cache.Client.SetKey(cacheKey, string(serialized), cacheFeedTTL)

	Send(w, string(serialized), http.StatusOK)
}

func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}
