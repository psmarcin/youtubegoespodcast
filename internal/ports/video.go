package ports

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/psmarcin/youtubegoespodcast/internal/app"
)

type videoDependencies interface {
	GetDetails(ctx context.Context, videoId string) (app.Details, error)
}

// videoHandler is server route handler for video redirection
func videoHandler(deps videoDependencies) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		videoID := ctx.Params("videoId")

		details, err := deps.GetDetails(ctx.Context(), videoID)
		if err != nil {
			l.WithError(err).Errorf("getting video url: %s", videoID)
			return ctx.SendStatus(http.StatusNotFound)
		}

		url := details.Url.String()

		if url == "" {
			l.Infof("didn't find video (%s) with audio", videoID)
			return ctx.SendStatus(http.StatusNotFound)
		}

		resp, err := http.Get(url)
		if err != nil {
			return fiber.NewError(http.StatusInternalServerError, err.Error())
		}

		return ctx.Redirect(resp.Request.URL.String(), http.StatusFound)
	}
}
