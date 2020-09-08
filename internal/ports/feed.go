package ports

import (
	"github.com/gofiber/fiber"
	"github.com/psmarcin/youtubegoespodcast/internal/app"
	"net/http"
)

const (
	ParamChannelId = "channelId"
)

// feedHandler is server route handler rss feed
func feedHandler(dependencies app.FeedService) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		channelID := ctx.Params(ParamChannelId)

		f, err := dependencies.Create(channelID)
		if err != nil {
			l.WithError(err).Errorf("can't create feed for %s", channelID)
			ctx.Next(err)
			return
		}

		response, err := f.Serialize()
		if err != nil {
			l.WithError(err).Errorf("can't serialize feed for %s", channelID)
			ctx.Next(err)
			return
		}

		ctx.Set("Content-Type", "application/rss+xml; charset=UTF-8")
		ctx.Status(http.StatusOK).SendBytes(response)
	}
}
