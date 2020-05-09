package api

import (
	"github.com/gofiber/fiber"
	"github.com/sirupsen/logrus"
	"ygp/pkg/youtube"
)

// TrendingHandler is a handler for GET /trending endpoint
func TrendingHandler(ctx *fiber.Ctx) {
	response, err := youtube.Yt.TrendingListFromCache()
	if err != nil {
		logrus.WithError(err).Errorf("can't get trending channels")
		ctx.Next(err)
	}

	err = ctx.JSON(response)
	if err != nil {
		ctx.Next(err)
	}
}
