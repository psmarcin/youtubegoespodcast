package adapters

import (
	"context"
	"github.com/pkg/errors"
	"github.com/psmarcin/youtubegoespodcast/internal/app"
	"google.golang.org/api/youtube/v3"
	"time"
)

const (
	YouTubeChannelsMaxResults = 3
	YouTubeChannelBaseURL     = "https://www.youtube.com/channel/"
)

type YouTubeAPIRepository struct {
	api *youtube.Service
}

func NewYouTube() (*youtube.Service, error) {
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx)
	if err != nil {
		l.WithError(err).Errorf("Can't create youtube service")
		return youtubeService, err
	}
	return youtubeService, nil
}

func NewYouTubeAPIRepository(yt *youtube.Service) YouTubeAPIRepository {
	return YouTubeAPIRepository{api: yt}
}

func (yt YouTubeAPIRepository) GetChannel(id string) (app.YouTubeChannel, error) {
	var channel app.YouTubeChannel

	call := yt.api.Channels.
		List([]string{"id", "snippet"}).
		MaxResults(1).
		Id(id)

	response, err := call.Do()
	if err != nil {
		l.WithError(err).Errorf("youtube api request failed")
		return channel, err
	}

	for _, item := range response.Items {
		channel, err = mapChannelItemToYouTubeChannels(item)
		if err != nil {
			return channel, err
		}
	}
	return channel, err
}

func (yt YouTubeAPIRepository) ListChannel(query string) ([]app.YouTubeChannel, error) {
	var channels []app.YouTubeChannel

	call := yt.api.Search.
		List([]string{"id", "snippet"}).
		MaxResults(YouTubeChannelsMaxResults).
		Type("channel").
		Q(query).
		Order("viewCount")

	response, err := call.Do()
	if err != nil {
		l.WithError(err).Errorf("youtube api request failed")
		return channels, err
	}

	for _, item := range response.Items {
		channel, err := mapSearchItemToYouTubeChannels(item)
		if err != nil {
			return channels, err
		}

		channels = append(channels, channel)
	}
	return channels, err
}

func mapChannelItemToYouTubeChannels(item *youtube.Channel) (app.YouTubeChannel, error) {
	ytChannel := app.YouTubeChannel{}
	thumbnail := getMaxThumbnailResolution(*item.Snippet.Thumbnails)
	publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
	if err != nil {
		return ytChannel, errors.Wrapf(err, "unable to parse: %s for channel: %s", item.Snippet.PublishedAt, item.Id)
	}
	ytChannel = app.YouTubeChannel{
		ChannelId:   item.Id,
		Country:     item.Snippet.Country,
		Description: item.Snippet.Description,
		PublishedAt: publishedAt,
		Thumbnail:   thumbnail.Url,
		Title:       item.Snippet.Title,
		Url:         YouTubeChannelBaseURL + item.Id,
		Author:      item.Snippet.CustomUrl,
		AuthorEmail: "email@example.com",
	}

	return ytChannel, nil
}

func mapSearchItemToYouTubeChannels(item *youtube.SearchResult) (app.YouTubeChannel, error) {
	ytChannel := app.YouTubeChannel{}
	thumbnail := getMaxThumbnailResolution(*item.Snippet.Thumbnails)
	publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
	if err != nil {
		return ytChannel, errors.Wrapf(err, "unable to parse: %s for channel: %s", item.Snippet.PublishedAt, item.Id)
	}
	ytChannel = app.YouTubeChannel{
		ChannelId:   item.Id.ChannelId,
		Description: item.Snippet.Description,
		PublishedAt: publishedAt,
		Thumbnail:   thumbnail.Url,
		Title:       item.Snippet.Title,
	}

	return ytChannel, nil
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
