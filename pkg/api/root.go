package api

import (
	"github.com/gofiber/fiber"
	"time"
)

type Root struct {
	Status    bool      `json:"status"`
	StartedAt time.Time `json:"startedAt"`
}

var r Root

func init() {
	r = Root{
		Status:    true,
		StartedAt: time.Now(),
	}
}

func rootHandler(ctx *fiber.Ctx) {
	err := ctx.JSON(r)
	if err != nil {
		ctx.Next(err)
	}
}
