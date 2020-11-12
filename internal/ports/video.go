package ports

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/psmarcin/youtubegoespodcast/internal/app"
	"net/http"
)

type videoDependencies interface {
	GetFileInformation(ctx context.Context, videoId string) (app.YTDLVideo, error)
}

// videoHandler is server route handler for video redirection
func videoHandler(deps videoDependencies) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		videoID := ctx.Params("videoId")

		details, err := deps.GetFileInformation(ctx.Context(), videoID)
		if err != nil {
			l.WithError(err).Errorf("getting video url: %s", videoID)
			return ctx.SendStatus(http.StatusNotFound)
		}

		url := details.FileUrl.String()

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
