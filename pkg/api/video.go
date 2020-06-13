package api

import (
	"github.com/gofiber/fiber"
	"net/http"

	"github.com/psmarcin/youtubegoespodcast/pkg/video"
)

// videoHandler is server route handler for video redirection
func videoHandler(dependencies video.Dependencies) func(ctx *fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		videoID := ctx.Params("videoId")

		details, err := video.GetDetails(videoID, dependencies)
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
