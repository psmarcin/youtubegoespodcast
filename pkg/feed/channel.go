package feed

import (
	"encoding/json"
	"errors"
	"github.com/eduncan911/podcast"
	"net/http"
	"time"

	"ygp/pkg/cache"
	"ygp/pkg/errx"
	"ygp/pkg/youtube"

	"github.com/sirupsen/logrus"
)

const (
	channelCachePRefix = "ytchannel_"
	channelCacheTTL    = time.Hour * 24 * 31
)

type ChannelDetailsResponse struct {
	Items []ChannelDetailsItems `json:"items"`
}

type ChannelDetailsHigh struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type ChannelDetailsThumbnails struct {
	High ChannelDetailsHigh `json:"high"`
}

type ChannelDetailsSnippet struct {
	Title       string                   `json:"title"`
	Description string                   `json:"description"`
	CustomURL   string                   `json:"customUrl"`
	PublishedAt time.Time                `json:"publishedAt"`
	Thumbnails  ChannelDetailsThumbnails `json:"thumbnails"`
	Country     string                   `json:"country"`
}

type ChannelDetailsItems struct {
	Snippet ChannelDetailsSnippet `json:"snippet"`
}

func (f *Feed) AddItem(item podcast.Item) error {
	if item.Title != "" && item.Enclosure.URL != "" {
		_, err := f.Content.AddItem(item)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Feed) GetDetails(channelID string) errx.APIError {
	channel := ChannelDetailsResponse{}
	_, err := cache.Client.GetKey(channelCachePRefix+channelID, &channel)
	// got cached value, fast return
	if err != nil {
		GetDetailsRequest(channelID, &channel)
	}

	if len(channel.Items) == 0 {
		return errx.New(errors.New("Can't find items for channel "+channelID), http.StatusNotFound)
	}

	item := channel.Items[0].Snippet
	fee := podcast.New(item.Title, ytChannelURL+f.ChannelID, item.Description, &item.PublishedAt, &item.PublishedAt)

	// TODO: map categories from youtube to apple categories
	fee.AddCategory("Arts", []string{"Design"})
	fee.Language = item.Country
	fee.Image = &podcast.Image{
		URL:         getImageURL(item.Thumbnails.High.URL),
		Title:       item.Title,
		Link:        ytChannelURL + f.ChannelID,
		Description: item.Description,
		Width:       0,
		Height:      0,
	}
	f.Content = fee

	return errx.APIError{}
}

func GetDetailsRequest(channelID string, channel *ChannelDetailsResponse) errx.APIError {
	req, err := http.NewRequest("GET", youtube.YouTubeURL+"channels", nil)
	if err != nil {
		logrus.WithError(err).Fatal("[YT] Can't create new request")
	}

	query := req.URL.Query()
	query.Add("part", "snippet")
	query.Add("id", channelID)
	query.Add("maxResults", "1")
	req.URL.RawQuery = query.Encode()

	requestError := youtube.Request(req, channel)
	if requestError.IsError() {
		return requestError
	}

	str, err := json.Marshal(channel)
	if err != nil {
		return errx.New(err, http.StatusInternalServerError)
	}

	go cache.Client.SetKey(channelCachePRefix+channelID, string(str), channelCacheTTL)

	if len(channel.Items) == 0 {
		return errx.New(errors.New("Can't find channel"), http.StatusInternalServerError)
	}

	return errx.APIError{}
}
