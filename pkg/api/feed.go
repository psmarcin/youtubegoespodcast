package api

import (
	"github.com/gofiber/fiber"
	"net/http"
	"time"
	"ygp/pkg/cache"
	"ygp/pkg/feed"
)

const (
	cacheFeedPrefix = "feed_"
	cacheFeedTTL    = time.Hour * 24
)

const (
	ParamChannelId = "channelId"
	QuerySearch    = "search"
)

func FeedHandler(ctx *fiber.Ctx) {
	channelID := ctx.Params(ParamChannelId)
	searchPhrase := ctx.FormValue(QuerySearch)
	cacheKey := cacheFeedPrefix + channelID + "_" + searchPhrase

	ctx.Set("Content-Type", "application/rss+xml; charset=UTF-8")

	// get cache
	c, _ := cache.Client.GetKey(cacheKey, nil)

	if c != "" {
		ctx.Status(http.StatusOK).Send(c)
		return
	}

	f := feed.New(channelID)
	err := f.GetDetails(channelID)

	if err != nil {
		ctx.Next(err)
		return
	}
	videos, getVideoErr := f.GetVideos(searchPhrase)
	if getVideoErr != nil {
		ctx.Next(getVideoErr)
		return
	}

	setVideosErr := f.SetVideos(videos)
	if setVideosErr != nil {
		ctx.Next(setVideosErr)
		return
	}

	f.SortItems()

	serialized, serializeErr := f.Serialize()
	if serializeErr != nil {
		ctx.Next(serializeErr)
		return
	}

	// set cache
	go cache.Client.SetKey(cacheKey, string(serialized), cacheFeedTTL)

	ctx.Status(http.StatusOK).SendBytes(serialized)
}
