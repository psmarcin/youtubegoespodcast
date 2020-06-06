package youtube

import (
	"github.com/psmarcin/youtubegoespodcast/pkg/cache"
	"time"
)

const (
	CacheChannelPrefix  = "youtube-channel-"
	CacheChannelTtl     = time.Hour * 24 * 31
	CacheChannelsPrefix = "youtube-channels-q-"
	CacheChannelsTtl    = time.Hour * 24
	CacheTrendingPrefix = "youtube-trending"
	CacheTrendingTtl    = time.Hour * 24
)

func (yt *YT) ChannelsListFromCache(query string) ([]Channel, error) {
	var channels []Channel
	cacheKey := CacheChannelsPrefix + query
	_, err := cache.Client.GetKey(cacheKey, &channels)
	if err == nil {
		return channels, nil
	}

	response, err := yt.ChannelsList(query)
	if err != nil {
		return channels, err
	}

	go cache.Client.MarshalAndSetKey(cacheKey, response, CacheChannelsTtl)

	return response, nil
}

func (yt *YT) ChannelsGetFromCache(id string) (Channel, error) {
	var channel Channel
	cacheKey := CacheChannelPrefix + id
	_, err := cache.Client.GetKey(cacheKey, &channel)
	if err == nil {
		return channel, nil
	}

	response, err := yt.ChannelGet(id)
	if err != nil {
		return channel, err
	}

	go cache.Client.MarshalAndSetKey(cacheKey, response, CacheChannelTtl)

	return response, nil
}

func (yt *YT) TrendingListFromCache() ([]Channel, error) {
	var channels []Channel
	cacheKey := CacheTrendingPrefix
	_, err := cache.Client.GetKey(cacheKey, &channels)
	if err == nil {
		return channels, nil
	}

	response, err := yt.TrendingList()
	if err != nil {
		return nil, err
	}

	go cache.Client.MarshalAndSetKey(cacheKey, response, CacheTrendingTtl)

	return response, nil
}
