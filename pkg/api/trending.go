package api

import (
	"github.com/gofiber/fiber"
	"net/http"
	"ygp/pkg/youtube"
)

// TrendingHandler is a handler for GET /trending endpoint
func TrendingHandler(ctx *fiber.Ctx) {
	response, err := youtube.GetTrending()
	if err != nil {
		ctx.Next(err)
		return
	}

	serialized := youtube.Serialize(response)
	ctx.Status(http.StatusOK).Send(serialized)
}
