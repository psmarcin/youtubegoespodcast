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

	details, err := video.GetDetails(videoID)
	if err != nil {
		logrus.Errorf("[API] error getting video url: %+v", err)
		ctx.SendStatus(http.StatusNotFound)
		return
	}

	url := details.FileUrl.String()

	if url == "" {
		logrus.Infof("didn't find video (%s) with audio", videoID)
		ctx.SendStatus(http.StatusNotFound)
		return
	}
	ctx.Redirect(url, http.StatusFound)
}
