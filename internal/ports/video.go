package ports

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/psmarcin/youtubegoespodcast/internal/app"
)

const (
	HTTPClientTimeout = 3 * time.Second
)

type videoDependencies interface {
	GetDetails(ctx context.Context, videoID string) (app.Details, error)
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

		url := details.URL.String()

		if url == "" {
			l.Infof("didn't find video (%s) with audio", videoID)
			return ctx.SendStatus(http.StatusNotFound)
		}
		client := http.Client{
			Timeout: HTTPClientTimeout,
		}
		resp, err := client.Get(url)
		if err != nil {
			return fiber.NewError(http.StatusInternalServerError, err.Error())
		}
		defer resp.Body.Close()

		return ctx.Redirect(resp.Request.URL.String(), http.StatusFound)
	}
}
