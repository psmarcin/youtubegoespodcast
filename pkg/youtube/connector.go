package youtube

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/youtube/v3"
	"strconv"
	"time"
)

const (
	ChannelsMaxResults    = 3
)

type YT struct {
	service *youtube.Service
}

// Channel holds metadata about channel
type Channel struct {
	ChannelId   string    `json:"channelId"`
	Country     string    `json:"country"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"publishedAt"`
	Thumbnail   string    `json:"thumbnail"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
}

// Video holds metadata about channel
type Video struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	PublishedAt time.Time         `json:"publishedAt"`
	Thumbnail   youtube.Thumbnail `json:"thumbnail"`
	Url         string            `json:"url"`
}

var Yt YT
var l = logrus.WithField("source", "youtube")

func init() {
	youtubeClient, err := New()
	if err != nil {
		l.WithError(err).Errorf("Can't initialize youtube service")
	}

	Yt = youtubeClient
}

func New() (YT, error) {
	yt := YT{}
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx)
	if err != nil {
		l.WithError(err).Errorf("Can't create youtube service")
		return yt, err
	}
	yt.service = youtubeService

	return yt, nil
}

func (yt *YT) ChannelGet(id string) (Channel, error) {
	var channel Channel

	l.WithField("id", id).Infof("channel get")

	call := yt.service.Channels.
		List("id,snippet").
		MaxResults(1).
		Id(id)

	response, err := call.Do()
	if err != nil {
		l.WithError(err).Errorf("youtube api request failed")
		return channel, err
	}

	for _, item := range response.Items {
		thumbnail := getMaxThumbnailResolution(*item.Snippet.Thumbnails)
		publishedAt, _ := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		channel = Channel{
			ChannelId:   item.Id,
			Country:     item.Snippet.Country,
			Description: item.Snippet.Description,
			PublishedAt: publishedAt,
			Thumbnail:   thumbnail.Url,
			Title:       item.Snippet.Title,
			Url:         item.Snippet.CustomUrl,
		}
	}
	return channel, err
}

func (yt *YT) ChannelsList(query string) ([]Channel, error) {
	var channels []Channel

	l.WithField("query", query).Infof("channels list")

	call := yt.service.Search.
		List("id,snippet").
		MaxResults(ChannelsMaxResults).
		Type("channel").
		Q(query).
		Order("viewCount")

	response, err := call.Do()
	if err != nil {
		l.WithError(err).Errorf("youtube api request failed")
		return channels, err
	}

	for _, item := range response.Items {
		thumbnail := getMaxThumbnailResolution(*item.Snippet.Thumbnails)
		publishedAt, _ := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		channels = append(channels, Channel{
			ChannelId:   item.Id.ChannelId,
			Description: item.Snippet.Description,
			PublishedAt: publishedAt,
			Thumbnail:   thumbnail.Url,
			Title:       item.Snippet.Title,
		})
	}
	return channels, err
}

func (yt *YT) TrendingList() ([]Channel, error) {
	var channels []Channel

	l.Infof("trending list")

	call := yt.service.Videos.
		List("id,snippet").
		MaxResults(ChannelsMaxResults).
		Chart("mostPopular")

	response, err := call.Do()
	if err != nil {
		l.WithError(err).Errorf("youtube api request failed")
		return channels, err
	}

	for _, item := range response.Items {
		thumbnail := getMaxThumbnailResolution(*item.Snippet.Thumbnails)
		publishedAt, _ := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		channels = append(channels, Channel{
			ChannelId:   item.Snippet.ChannelId,
			Thumbnail:   thumbnail.Url,
			PublishedAt: publishedAt,
			Title:       item.Snippet.ChannelTitle,
		})
	}
	return channels, err
}

func (yt *YT) VideosList(channelId string) ([]Video, error) {
	var videos []Video

	f,err := GetFeed(channelId)
	if err != nil {
		l.WithError(err).Errorf("can't get video list for %s", channelId)
		return videos, err
	}

	for _, item := range f.Entry {
		publishedAt, _ := time.Parse(time.RFC3339, item.Published)

		tHeight, err := strconv.Atoi(item.Group.Thumbnail.Height)
		if err != nil{
			l.WithError(err).Errorf("can't parse video %s thumbnail height %s", item.ID, item.Group.Thumbnail.Height)
		}

		tWidth, err := strconv.Atoi(item.Group.Thumbnail.Width)
		if err != nil{
			l.WithError(err).Errorf("can't parse video %s thumbnail width %s", item.ID, item.Group.Thumbnail.Width)
		}

		videos = append(videos, Video{
			ID:          item.VideoId,
			Title:       item.Title,
			Description: item.Group.Description,
			PublishedAt: publishedAt,
			Thumbnail:   youtube.Thumbnail{
				Height:          int64(tHeight),
				Url:             item.Group.Thumbnail.URL,
				Width:           int64(tWidth),
			},
			Url:         item.Link.Href,
		})
	}
	return videos, err
}

func getMaxThumbnailResolution(thumbnails youtube.ThumbnailDetails) youtube.Thumbnail {
	if thumbnails.Maxres != nil {
		return *thumbnails.Maxres
	}
	if thumbnails.High != nil {
		return *thumbnails.High
	}
	if thumbnails.Medium != nil {
		return *thumbnails.Medium
	}
	if thumbnails.Standard != nil {
		return *thumbnails.Standard
	}
	return youtube.Thumbnail{}
}
