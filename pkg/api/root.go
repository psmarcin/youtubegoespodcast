package api

import (
	"github.com/gofiber/fiber"
	"github.com/psmarcin/youtubegoespodcast/pkg/youtube"
)

const (
	BaseFeedURL = "https://yt.psmarcin.dev/feed/channel/"
)

// rootHandler is server route handler for main page and interaction with user
func rootHandler(ctx *fiber.Ctx) {
	var err error
	var channels []youtube.Channel

	channelId := ctx.FormValue("channelId")
	if channelId != "" {
		channelId = BaseFeedURL + channelId
	}

	q := ctx.FormValue("q")
	if q != "" {
		channels, err = youtube.Yt.ChannelsListFromCache(q)
	}
	if err != nil {
		ctx.Next(err)
		return
	}

	ctx.Set("content-type", "text/html; charset=utf-8")
	err = ctx.Render("index", fiber.Map{
		"Channels":  channels,
		"ChannelId": channelId,
	})

	if err != nil {
		l.WithError(err).Errorf("error while rendering template")
		ctx.Next(err)
	}
}
