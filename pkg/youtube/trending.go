package youtube

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"ygp/pkg/cache"
)

const (
	trendingCacheTTL           = time.Hour * 24
	trendingCachePrefix string = "yttrending_main"
)

// GetTrending collects trending from YouTube API
func GetTrending() (YoutubeResponse, error) {
	trending := YoutubeResponse{}
	_, err := cache.Client.GetKey(trendingCachePrefix, &trending)
	if err == nil {
		return trending, nil
	}

	req, err := http.NewRequest("GET", YouTubeURL+"videos", nil)
	if err != nil {
		logrus.WithError(err).Fatal("[YT] Can't create new request")
		return trending, err
	}
	query := req.URL.Query()
	query.Add("part", "snippet")
	query.Add("chart", "mostPopular")
	req.URL.RawQuery = query.Encode()

	requestErr := Request(req, &trending)

	if requestErr != nil {
		return trending, requestErr
	}

	// save videoDetails to cache
	str, err := json.Marshal(trending)
	if err != nil {
		logrus.WithError(err).Fatal("[YT] Can't create new request")
		return trending, err
	}
	go cache.Client.SetKey(trendingCachePrefix, string(str), trendingCacheTTL)

	return trending, nil
}
