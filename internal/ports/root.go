package ports

import (
	"github.com/gofiber/fiber"
	"github.com/psmarcin/youtubegoespodcast/internal/app"
)

const (
	BaseFeedURL = "https://yt.psmarcin.dev/feed/channel/"
)

type rootDependencies interface {
	ListChannel(query string) ([]app.YouTubeChannel, error)
}

// rootHandler is server route handler for main page and interaction with user
func rootHandler(rootDependency rootDependencies) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		var err error
		var channels []app.YouTubeChannel

		channelId := ctx.FormValue("channelId")
		if channelId != "" {
			channelId = BaseFeedURL + channelId
		}

		q := ctx.FormValue("q")
		if q != "" {
			channels, err = rootDependency.ListChannel(q)
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
}
