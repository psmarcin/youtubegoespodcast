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
	cacheFeedTTL    = time.Hour * 24 * 1
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

	if err.IsError() {
		ctx.Next(err.Err)
		return
	}
	videos, getVideoErr := f.GetVideos(searchPhrase, false)
	if getVideoErr.IsError() {
		ctx.Next(getVideoErr.Err)
		return
	}

	setVideosErr := f.SetVideos(videos)
	if setVideosErr.IsError() {
		ctx.Next(setVideosErr.Err)
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
