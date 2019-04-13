package youtube

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type YoutubeResponse struct {
	Items []Items `json:"items"`
}
type High struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
type Thumbnails struct {
	High High `json:"high"`
}
type Snippet struct {
	ChannelID    string     `json:"channelId"`
	Title        string     `json:"title"`
	Thumbnails   Thumbnails `json:"thumbnails"`
	ChannelTitle string     `json:"channelTitle"`
}
type Items struct {
	Snippet Snippet `json:"snippet"`
}

type Channel struct {
	ChannelId string `json:"channelId"`
	Thumbnail string `json:"thumbnail"`
	Title     string `json:"title"`
}

// GetTrendings collects trendings from YouTube API
func GetTrendings() YoutubeResponse {
	trendings := YoutubeResponse{}
	req, err := http.NewRequest("GET", YouTubeURL+"videos", nil)
	if err != nil {
		logrus.WithError(err).Fatal("[YT] Can't create new request")
	}
	query := req.URL.Query()
	query.Add("part", "snippet")
	query.Add("chart", "mostPopular")
	req.URL.RawQuery = query.Encode()

	Request(req, &trendings)

	return trendings
}

// SerializeTrending returns mapped channels with channelId, thumbnail and title
func SerializeTrending(channelResponse YoutubeResponse) string {
	channels := []Channel{}
	for _, item := range channelResponse.Items {
		channels = append(channels, Channel{
			ChannelId: item.Snippet.ChannelID,
			Thumbnail: item.Snippet.Thumbnails.High.URL,
			Title:     item.Snippet.ChannelTitle,
		})
	}
	serialized, err := json.Marshal(channels)
	if err != nil {
		logrus.WithError(err).Info("[YT] Serialization")
	}
	return string(serialized)
}
