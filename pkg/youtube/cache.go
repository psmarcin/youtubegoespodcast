package youtube

import (
	"github.com/sirupsen/logrus"
	"time"
	"ygp/pkg/cache"
)

const (
	CacheChannelsPrefix = "youtube-channels-q-"
	CacheChannelsTtl = time.Hour * 24
	CacheTrendingPrefix = "youtube-trending"
	CacheTrendingTtl = time.Hour * 24
)

func (yt *YT) ChannelsListFromCache(query string) ([]Channel, error){
	var channels []Channel
	cacheKey := CacheChannelsPrefix+query
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


func (yt *YT) TrendingListFromCache() ([]Channel, error){
	var channels []Channel
	cacheKey := CacheTrendingPrefix
	_, err := cache.Client.GetKey(cacheKey, &channels)
	if err == nil {
		return channels, nil
	}

	response, err := yt.TrendingList()
	logrus.WithField("respose", response).Infof("response from trending")
	if err != nil {
		return nil, err
	}

	go cache.Client.MarshalAndSetKey(cacheKey, response, CacheTrendingTtl)

	return response, nil
}

