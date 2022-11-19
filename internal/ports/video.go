package ports

import (
	"context"
	"net/http"

	fib "github.com/gofiber/fiber/v2"
	"github.com/psmarcin/youtubegoespodcast/internal/app"
)

type videoDependencies interface {
	GetDetails(context.Context, string) (app.Details, error)
}

// videoHandler is server route handler for video redirection
func videoHandler(deps videoDependencies) func(ctx *fib.Ctx) error {
	return func(ctx *fib.Ctx) error {
		videoID := ctx.Params("videoId")

		details, err := deps.GetDetails(ctx.Context(), videoID)
		if err != nil {
			l.WithError(err).Errorf("getting video url: %s", videoID)
			return ctx.SendStatus(http.StatusNotFound)
		}

		url := details.URL.String()

		if url == "" {
			l.Infof("didn't find video (%s) with audio", videoID)
			return ctx.SendStatus(http.StatusNotFound)
		}

		resp, err := http.Get(url) //nolint
		if err != nil {
			return fib.NewError(http.StatusInternalServerError, err.Error())
		}

		return ctx.Redirect(resp.Request.URL.String(), http.StatusFound)
	}
}
