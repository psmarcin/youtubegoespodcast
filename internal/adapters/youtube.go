package adapters

import (
	"context"
	"github.com/pkg/errors"
	"github.com/psmarcin/youtubegoespodcast/internal/app"
	feedDomain "github.com/psmarcin/youtubegoespodcast/internal/domain/feed"
	"go.opentelemetry.io/otel/label"
	"google.golang.org/api/youtube/v3"
	"net/url"
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

func (yt YouTubeAPIRepository) GetChannel(ctx context.Context, id string) (app.YouTubeChannel, error) {
	var channel app.YouTubeChannel
	ctx, span := tracer.Start(ctx, "get-channel")
	span.SetAttributes(label.String("id", id))
	defer span.End()

	call := yt.api.Channels.
		List([]string{"id", "snippet", "topicDetails"}).
		MaxResults(1).
		Id(id)

	response, err := call.Do()
	if err != nil {
		l.WithError(err).Errorf("youtube api request failed")
		span.RecordError(err)
		return channel, err
	}

	for _, item := range response.Items {
		channel, err = mapChannelToYouTubeChannel(item)
		if err != nil {
			span.RecordError(err)
			return channel, err
		}
	}
	return channel, err
}

func (yt YouTubeAPIRepository) ListChannel(ctx context.Context, query string) ([]app.YouTubeChannel, error) {
	var channels []app.YouTubeChannel
	ctx, span := tracer.Start(ctx, "list-channel")
	span.SetAttributes(label.String("query", query))
	defer span.End()

	call := yt.api.Search.
		List([]string{"id", "snippet"}).
		MaxResults(YouTubeChannelsMaxResults).
		Type("channel").
		Q(query).
		Order("viewCount")

	response, err := call.Do()
	if err != nil {
		l.WithError(err).Errorf("youtube api request failed")
		span.RecordError(err)
		return channels, err
	}

	for _, item := range response.Items {
		channel, err := mapSearchItemToYouTubeChannel(item)
		if err != nil {
			span.RecordError(err)
			return channels, err
		}

		channels = append(channels, channel)
	}
	return channels, err
}

func mapChannelToYouTubeChannel(item *youtube.Channel) (app.YouTubeChannel, error) {
	ytChannel := app.YouTubeChannel{}

	thumbnail := getMaxThumbnailResolution(*item.Snippet.Thumbnails)
	thumbnailUrl, err := url.Parse(thumbnail.Url)
	if err != nil {
		l.WithError(err).Errorf("can't parse thumbnail url: %s", thumbnail.Url)
	}

	publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
	if err != nil {
		return ytChannel, errors.Wrapf(err, "unable to parse: %s for channel: %s", item.Snippet.PublishedAt, item.Id)
	}

	category := feedDomain.SelectCategory(item.TopicDetails.TopicCategories)
	l.WithField("category", category).Infof("found category for channel: %s", item.Id)

	ytChannel = app.YouTubeChannel{
		ChannelId:   item.Id,
		Country:     item.Snippet.Country,
		Description: item.Snippet.Description,
		PublishedAt: publishedAt,
		Thumbnail: app.YouTubeThumbnail{
			Height: int(thumbnail.Height),
			Width:  int(thumbnail.Width),
			Url:    *thumbnailUrl,
		},
		Title:       item.Snippet.Title,
		Url:         YouTubeChannelBaseURL + item.Id,
		Author:      item.Snippet.CustomUrl,
		AuthorEmail: "email@example.com",
		Category:    category,
	}

	return ytChannel, nil
}

func mapSearchItemToYouTubeChannel(item *youtube.SearchResult) (app.YouTubeChannel, error) {
	ytChannel := app.YouTubeChannel{}

	thumbnail := getMaxThumbnailResolution(*item.Snippet.Thumbnails)
	thumbnailUrl, err := url.Parse(thumbnail.Url)
	if err != nil {
		l.WithError(err).Errorf("can't parse thumbnail url: %s", thumbnail.Url)
	}

	publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
	if err != nil {
		return ytChannel, errors.Wrapf(err, "unable to parse: %s for channel: %s", item.Snippet.PublishedAt, item.Id)
	}

	category := feedDomain.SelectCategory([]string{})

	ytChannel = app.YouTubeChannel{
		ChannelId:   item.Id.ChannelId,
		Description: item.Snippet.Description,
		PublishedAt: publishedAt,
		Thumbnail: app.YouTubeThumbnail{
			Height: int(thumbnail.Height),
			Width:  int(thumbnail.Width),
			Url:    *thumbnailUrl,
		},
		Title:    item.Snippet.Title,
		Category: category,
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
	if thumbnails.Default != nil {
		return *thumbnails.Default
	}
	return youtube.Thumbnail{}
}
