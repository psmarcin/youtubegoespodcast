package api

import (
	"github.com/gofiber/fiber"
	"github.com/psmarcin/youtubegoespodcast/internal/app"
	"net/http"
)

type videoDependencies interface {
	GetFileInformation(videoId string) (app.YTDLVideo, error)
}

// videoHandler is server route handler for video redirection
func videoHandler(deps videoDependencies) func(ctx *fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		videoID := ctx.Params("videoId")

		details, err := deps.GetFileInformation(videoID)
		if err != nil {
			l.WithError(err).Errorf("getting video url: %s", videoID)
			ctx.SendStatus(http.StatusNotFound)
			return
		}

		url := details.FileUrl.String()

		if url == "" {
			l.Infof("didn't find video (%s) with audio", videoID)
			ctx.SendStatus(http.StatusNotFound)
			return
		}

		ctx.Redirect(url, http.StatusFound)
	}
}
