package api

import (
	"github.com/gofiber/fiber"
	"ygp/pkg/youtube"
)

var disableCache = false

// TrendingHandler is a handler for GET /trending endpoint
func TrendingHandler(ctx *fiber.Ctx) {
	response, err := youtube.GetTrending(disableCache)
	if err.IsError() {
		ctx.Next(err.Err)
		return
	}

	jsonError := ctx.JSON(response)
	if jsonError != nil {
		ctx.Next(jsonError)
	}
}
