package api

import (
	"github.com/gofiber/fiber"
	"net/http"
	"ygp/pkg/youtube"
)

// Handler is default router handler for GET /channel endpoint
func ChannelsHandler(ctx *fiber.Ctx) {
	q := ctx.FormValue("q")
	channels, err := youtube.GetChannels(q)
	if err != nil {
		ctx.Next(err)
		return
	}

	serialized := youtube.Serialize(channels)
	ctx.Status(http.StatusOK).Send(serialized)
}
