package api

import (
	"github.com/gofiber/fiber"
	"net/http"
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

	serialized := youtube.Serialize(response)
	ctx.Status(http.StatusOK).Send(serialized)
}
