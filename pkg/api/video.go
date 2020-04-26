package api

import (
	"github.com/gofiber/fiber"
	"net/http"

	"github.com/sirupsen/logrus"

	"ygp/pkg/video"
)

// Handler handles GET /video?videoId= route
func VideoHandler(ctx *fiber.Ctx) {
	videoID := ctx.Params("videoId")

	videoURL, err := video.GetURL(videoID)
	if err != nil {
		logrus.Errorf("[API] error getting video url: %+v", err)
		ctx.SendStatus(http.StatusNotFound)
		return
	}

	if videoURL == "" {
		logrus.Infof("didn't find video (%s) with audio", videoID)
		ctx.SendStatus(http.StatusNotFound)
		return
	}
	ctx.Redirect(videoURL, http.StatusFound)
}
