package app

import (
	"context"
	"net/url"
	"time"

	"github.com/cockroachdb/errors"
	"go.opentelemetry.io/otel/attribute"
)

const (
	CacheGetChannelPrefix  = "youtube-channel-"
	CacheListChannelPrefix = "youtube-channels-q-"
)

type youTubeRepository interface {
	ListEntry(context.Context, string) ([]YouTubeFeedEntry, error)
}

type youTubeAPIRepository interface {
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
	URL    url.URL
}

// YouTubeChannel holds metadata about channel
type YouTubeChannel struct {
	Author      string
	AuthorEmail string
	ChannelID   string
	Country     string
	Description string
	PublishedAt time.Time
	Thumbnail   YouTubeThumbnail
	Title       string
	URL         string
	Category    string
}

type YouTubeService struct {
	youTubeRepository    youTubeRepository
	youTubeAPIRepository youTubeAPIRepository
	cache                CacheService
}

func NewYouTubeService(ytRepo youTubeRepository, ytAPIRepo youTubeAPIRepository, cache CacheService) YouTubeService {
	return YouTubeService{
		youTubeRepository:    ytRepo,
		youTubeAPIRepository: ytAPIRepo,
		cache:                cache,
	}
}

func (y YouTubeService) ListEntry(ctx context.Context, channelID string) ([]YouTubeFeedEntry, error) {
	return y.youTubeRepository.ListEntry(ctx, channelID)
}

func (y YouTubeService) GetChannel(ctx context.Context, channelID string) (YouTubeChannel, error) {
	return y.youTubeAPIRepository.GetChannel(ctx, channelID)
}

func (y YouTubeService) ListChannel(ctx context.Context, query string) ([]YouTubeChannel, error) {
	return y.youTubeAPIRepository.ListChannel(ctx, query)
}

func (y YouTubeService) GetChannelCache(ctx context.Context, channelID string) (YouTubeChannel, error) {
	ctx, span := tracer.Start(ctx, "youtube-get-channel-cache")
	span.SetAttributes(attribute.Any("channelID", channelID))
	defer span.End()

	var channel YouTubeChannel
	cacheKey := CacheGetChannelPrefix + channelID
	err := y.cache.Get(ctx, cacheKey, &channel)
	if err == nil {
		return channel, nil
	}

	response, err := y.GetChannel(ctx, channelID)
	if err != nil {
		span.RecordError(err)
		return channel, errors.Wrapf(err, "unable value get channel for channel id: %s", channelID)
	}

	go y.cache.MarshalAndSet(ctx, cacheKey, response)

	return response, nil
}

func (y YouTubeService) ListChannelCache(ctx context.Context, query string) ([]YouTubeChannel, error) {
	ctx, span := tracer.Start(ctx, "youtube-list-channel-chache")
	span.SetAttributes(attribute.Any("query", query))
	defer span.End()

	var channels []YouTubeChannel
	cacheKey := CacheListChannelPrefix + query

	err := y.cache.Get(ctx, cacheKey, &channels)
	if err == nil {
		return channels, nil
	}

	response, err := y.ListChannel(ctx, query)
	if err != nil {
		span.RecordError(err)
		return channels, errors.Wrapf(err, "unable value list channel for query: %s", query)
	}

	go y.cache.MarshalAndSet(ctx, cacheKey, response)

	return response, nil
}
