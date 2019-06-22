package feed

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"ytg/pkg/errx"
	"ytg/pkg/redis_client"
	"ytg/pkg/youtube"

	"github.com/sirupsen/logrus"
)

const (
	channelCachePRefix = "ytchannel_"
	channelCacheTTL    = time.Hour * 1
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

func (f *Feed) addItem(item Item) error {
	if item.Title != "" && item.Enclosure.URL != "" {
		f.Items = append(f.Items, item)
	}
	return nil
}

func (f *Feed) getDetails(channelID string) errx.APIError {
	channel := ChannelDetailsResponse{}
	_, err := redis_client.Client.GetKey(channelCachePRefix+channelID, &channel)
	// got cached value, fast return
	if err != nil {
		getDetailsRequest(channelID, &channel)
	}

	if len(channel.Items) == 0 {
		return errx.NewAPIError(errors.New("Can't find items for channel "+channelID), http.StatusNotFound)
	}

	item := channel.Items[0].Snippet

	f.Title = item.Title
	f.Link = ytChannelURL + f.ChannelID
	f.Description = item.Description
	f.Category = "category"
	f.Language = "en"
	f.LastBuildDate = item.PublishedAt.String()
	f.PubDate = item.PublishedAt.String()
	f.Image = Image{
		URL:   getImageURL(item.Thumbnails.High.URL),
		Title: item.Title,
		Link:  ytChannelURL + f.ChannelID,
	}
	f.ITAuthor = item.CustomURL
	f.ITSubtitle = item.Title
	f.ITSummary = ITSummary{
		Text: item.Description,
	}
	f.ITImage = ITImage{
		Href: getImageURL(item.Thumbnails.High.URL),
	}
	f.ITExplicit = "no"
	return errx.APIError{}
}

func getDetailsRequest(channelID string, channel *ChannelDetailsResponse) errx.APIError {
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
		return errx.NewAPIError(err, http.StatusInternalServerError)
	}
	go redis_client.Client.SetKey(channelCachePRefix+channelID, string(str), channelCacheTTL)

	if len(channel.Items) == 0 {
		return errx.NewAPIError(errors.New("Can't find channel"), http.StatusInternalServerError)
	}

	return errx.APIError{}
}
