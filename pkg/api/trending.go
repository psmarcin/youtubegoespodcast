package api

import (
	"github.com/gofiber/fiber"
	"ygp/pkg/youtube"
)

// TrendingHandler is a handler for GET /trending endpoint
func TrendingHandler(ctx *fiber.Ctx) {
	response, err := youtube.Yt.TrendingListFromCache()
	if err != nil {
		l.WithError(err).Errorf("can't get trending channels")
		ctx.Next(err)
	}

	err = ctx.JSON(response)
	if err != nil {
		ctx.Next(err)
	}
}
