package ports

import (
	"net/http"

	fib "github.com/gofiber/fiber/v2"
	fiber_opentelemetry "github.com/psmarcin/fiber-opentelemetry/pkg/fiber-otel"
	"github.com/psmarcin/youtubegoespodcast/internal/app"
)

const (
	ParamChannelID = "channelId"
)

// feedHandler is server route handler rss feed
func feedHandler(dependencies app.FeedService) func(*fib.Ctx) error {
	return func(ctx *fib.Ctx) error {
		channelID := ctx.Params(ParamChannelID)

		c := fiber_opentelemetry.FromCtx(ctx)
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
