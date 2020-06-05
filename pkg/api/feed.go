package api

import (
	"github.com/gofiber/fiber"
	"net/http"
	"ygp/pkg/feed"
	"ygp/pkg/youtube"
)

const (
	ParamChannelId = "channelId"
)

func FeedHandler(ctx *fiber.Ctx) {
	channelID := ctx.Params(ParamChannelId)

	f := feed.New(channelID)
	err := f.GetDetails(channelID)
	if err != nil {
		ctx.Next(err)
		return
	}

	videos, getVideoErr := youtube.Yt.VideosList(f.ChannelID)
	if getVideoErr != nil {
		ctx.Next(getVideoErr)
		return
	}

	setVideosErr := f.SetVideos(videos)
	if setVideosErr != nil {
		ctx.Next(setVideosErr)
		return
	}

	f.SortVideos()

	serialized, serializeErr := f.Serialize()
	if serializeErr != nil {
		ctx.Next(serializeErr)
		return
	}

	ctx.Set("Content-Type", "application/rss+xml; charset=UTF-8")
	ctx.Status(http.StatusOK).SendBytes(serialized)
}
