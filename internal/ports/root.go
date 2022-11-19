package ports

import (
	"context"

	fib "github.com/gofiber/fiber/v2"
	"github.com/psmarcin/youtubegoespodcast/internal/app"
)

const (
	BaseFeedURL = "https://yt.psmarcin.dev/feed/channel/"
)

type rootDependencies interface {
	ListChannel(ctx context.Context, query string) ([]app.YouTubeChannel, error)
}

// rootHandler is server route handler for main page and interaction with user
func rootHandler(rootDependency rootDependencies) func(*fib.Ctx) error {
	return func(ctx *fib.Ctx) error {
		var err error
		var channels []app.YouTubeChannel

		channelID := ctx.FormValue("channelId")
		if channelID != "" {
			channelID = BaseFeedURL + channelID
		}

		q := ctx.FormValue("q")
		if q != "" {
			channels, err = rootDependency.ListChannel(ctx.Context(), q)
		}
		if err != nil {
			return err
		}

		ctx.Set("content-type", "text/html; charset=utf-8")
		err = ctx.Render("templates/index", fib.Map{
			"Channels":  channels,
			"ChannelId": channelID,
		})

		if err != nil {
			l.WithError(err).Errorf("error while rendering template")
			return err
		}

		return nil
	}
}
