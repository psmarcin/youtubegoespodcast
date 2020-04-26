package youtube

import (
	"net/http"

	"ygp/pkg/errx"

	"github.com/sirupsen/logrus"
)

// GetChannels collects trending from YouTube API
func GetChannels(q string) (YoutubeResponse, errx.APIError) {
	trending := YoutubeResponse{}
	req, err := http.NewRequest("GET", YouTubeURL+"search", nil)
	if err != nil {
		logrus.WithError(err).Fatal("[YT] Can't create new request")
		return trending, errx.New(err, http.StatusBadRequest)
	}
	query := req.URL.Query()
	query.Add("part", "snippet")
	query.Add("type", "channel")
	query.Add("q", q)
	query.Add("fields", "items(snippet(channelId,channelTitle,thumbnails/high,title))")
	query.Add("order", "viewCount")
	query.Add("maxResults", "5")
	req.URL.RawQuery = query.Encode()

	requestErr := Request(req, &trending)
	if requestErr.IsError() {
		return trending, requestErr
	}

	return trending, errx.APIError{}
}
