package api

import (
	"github.com/gofiber/fiber"
	"github.com/rylio/ytdl"
	"net/http"

	"github.com/psmarcin/youtubegoespodcast/pkg/video"
)

// videoHandler is server route handler for video redirection
func videoHandler() func(ctx *fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		videoID := ctx.Params("videoId")
		v := video.New(videoID)

		details, err := v.GetFileInformation(ytdl.DefaultClient, ytdl.DefaultClient)
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
