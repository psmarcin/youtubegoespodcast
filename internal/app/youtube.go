package app

import (
	"context"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/label"
	"net/url"
	"time"
)

const (
	CacheGetChannelPrefix  = "youtube-channel-"
	CacheGetChannelTtl     = time.Hour * 24 * 31
	CacheListChannelPrefix = "youtube-channels-q-"
	CacheListChannelTtl    = time.Hour * 24
)

type youTubeRepository interface {
	ListEntry(context.Context, string) ([]YouTubeFeedEntry, error)
}

type youTubeApiRepository interface {
	GetChannel(context.Context, string) (YouTubeChannel, error)
	ListChannel(context.Context, string) ([]YouTubeChannel, error)
}

type YouTubeFeedEntry struct {
	ID          string
	Title       string
	Description string
	Published   time.Time
	URL         url.URL
	Thumbnail   YouTubeThumbnail
}

type YouTubeThumbnail struct {
	Height int
	Width  int
	Url    url.URL
}

// YouTubeChannel holds metadata about channel
type YouTubeChannel struct {
	Author      string
	AuthorEmail string
	ChannelId   string
	Country     string
	Description string
	PublishedAt time.Time
	Thumbnail   YouTubeThumbnail
	Title       string
	Url         string
	Category    string
}

type YouTubeService struct {
	youTubeRepository    youTubeRepository
	youTubeApiRepository youTubeApiRepository
	cache                CacheService
}

func NewYouTubeService(ytRepo youTubeRepository, ytApiRepo youTubeApiRepository, cache CacheService) YouTubeService {
	return YouTubeService{
		youTubeRepository:    ytRepo,
		youTubeApiRepository: ytApiRepo,
		cache:                cache,
	}
}

func (y YouTubeService) ListEntry(ctx context.Context, channelId string) ([]YouTubeFeedEntry, error) {
	return y.youTubeRepository.ListEntry(ctx, channelId)
}

func (y YouTubeService) GetChannel(ctx context.Context, channelId string) (YouTubeChannel, error) {
	return y.youTubeApiRepository.GetChannel(ctx, channelId)
}

func (y YouTubeService) ListChannel(ctx context.Context, query string) ([]YouTubeChannel, error) {
	return y.youTubeApiRepository.ListChannel(ctx, query)
}

func (y YouTubeService) GetChannelCache(ctx context.Context, channelId string) (YouTubeChannel, error) {
	ctx, span := tracer.Start(ctx, "youtube-get-channel-cache")
	span.SetAttributes(label.Any("channelId", channelId))
	defer span.End()

	var channel YouTubeChannel
	cacheKey := CacheGetChannelPrefix + channelId
	err := y.cache.Get(ctx, cacheKey, &channel)
	if err == nil {
		return channel, nil
	}

	response, err := y.GetChannel(ctx, channelId)
	if err != nil {
		span.RecordError(ctx, err)
		return channel, errors.Wrapf(err, "unable value get channel for channel id: %s", channelId)
	}

	go y.cache.MarshalAndSet(ctx, cacheKey, response)

	return response, nil
}

func (y YouTubeService) ListChannelCache(ctx context.Context, query string) ([]YouTubeChannel, error) {
	ctx, span := tracer.Start(ctx, "youtube-list-channel-chache")
	span.SetAttributes(label.Any("query", query))
	defer span.End()

	var channels []YouTubeChannel
	cacheKey := CacheListChannelPrefix + query

	err := y.cache.Get(ctx, cacheKey, &channels)
	if err == nil {
		return channels, nil
	}

	response, err := y.ListChannel(ctx, query)
	if err != nil {
		span.RecordError(ctx, err)
		return channels, errors.Wrapf(err, "unable value list channel for query: %s", query)
	}

	go y.cache.MarshalAndSet(ctx, cacheKey, response)

	return response, nil
}
