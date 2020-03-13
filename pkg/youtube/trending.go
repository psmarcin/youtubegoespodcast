package youtube

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"ygp/pkg/cache"
	"ygp/pkg/errx"
)

const (
	trendingCacheTTL    time.Duration = time.Hour * 24 * 30
	trendingCachePrefix               = "yttrending_main"
)

// GetTrending collects trending from YouTube API
func GetTrending() (YoutubeResponse, errx.APIError) {
	trending := YoutubeResponse{}
	_, err := cache.Client.GetKey(trendingCachePrefix, &trending)
	if err == nil {
		return trending, errx.APIError{}
	}

	req, err := http.NewRequest("GET", YouTubeURL+"videos", nil)
	if err != nil {
		logrus.WithError(err).Fatal("[YT] Can't create new request")
		return trending, errx.NewAPIError(err, http.StatusInternalServerError)
	}
	query := req.URL.Query()
	query.Add("part", "snippet")
	query.Add("chart", "mostPopular")
	req.URL.RawQuery = query.Encode()

	requestErr := Request(req, &trending)

	if requestErr.IsError() {
		return trending, requestErr
	}

	// save videoDetails to cache
	str, err := json.Marshal(trending)
	if err != nil {
		logrus.WithError(err).Fatal("[YT] Can't create new request")
		return trending, errx.NewAPIError(err, http.StatusInternalServerError)
	}
	go cache.Client.SetKey(trendingCachePrefix, string(str), trendingCacheTTL)
	return trending, errx.APIError{}
}
