package app

import (
	"github.com/pkg/errors"
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
	ListEntry(string) ([]YouTubeFeedEntry, error)
}

type youTubeApiRepository interface {
	GetChannel(string) (YouTubeChannel, error)
	ListChannel(string) ([]YouTubeChannel, error)
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

func (y YouTubeService) ListEntry(channelId string) ([]YouTubeFeedEntry, error) {
	return y.youTubeRepository.ListEntry(channelId)
}

func (y YouTubeService) GetChannel(channelId string) (YouTubeChannel, error) {
	return y.youTubeApiRepository.GetChannel(channelId)
}

func (y YouTubeService) ListChannel(query string) ([]YouTubeChannel, error) {
	return y.youTubeApiRepository.ListChannel(query)
}

func (y YouTubeService) GetChannelCache(channelId string) (YouTubeChannel, error) {
	var channel YouTubeChannel
	cacheKey := CacheGetChannelPrefix + channelId
	err := y.cache.Get(cacheKey, &channel)
	if err == nil {
		return channel, nil
	}

	response, err := y.GetChannel(channelId)
	if err != nil {
		return channel, errors.Wrapf(err, "unable value get channel for channel id: %s", channelId)
	}

	go y.cache.MarshalAndSet(cacheKey, response)

	return response, nil
}

func (y YouTubeService) ListChannelCache(query string) ([]YouTubeChannel, error) {
	var channels []YouTubeChannel
	cacheKey := CacheListChannelPrefix + query

	err := y.cache.Get(cacheKey, &channels)
	if err == nil {
		return channels, nil
	}

	response, err := y.ListChannel(query)
	if err != nil {
		return channels, errors.Wrapf(err, "unable value list channel for query: %s", query)
	}

	go y.cache.MarshalAndSet(cacheKey, response)

	return response, nil
}
