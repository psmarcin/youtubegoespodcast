package ports

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/psmarcin/youtubegoespodcast/internal/app"
	"net/http"
)

const (
	ParamChannelId = "channelId"
)

// feedHandler is server route handler rss feed
func feedHandler(dependencies app.FeedService) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		channelID := ctx.Params(ParamChannelId)

		c := ctx.Locals("ctx").(context.Context)
		f, err := dependencies.Create(c, channelID)
		if err != nil {
			l.WithError(err).Errorf("can't create feed for %s", channelID)
			return err
		}

		response, err := f.Serialize()
		if err != nil {
			l.WithError(err).Errorf("can't serialize feed for %s", channelID)
			return err
		}

		ctx.Set("Content-Type", "application/rss+xml; charset=UTF-8")
		return ctx.Status(http.StatusOK).Send(response)
	}
}
