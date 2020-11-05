package ports

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
)

func RequestContext() func(c *fiber.Ctx) error {
	tracer := global.TracerProvider().Tracer(
		"yt.psmarcin.dev/api",
		trace.WithInstrumentationVersion("1.0.0"),
	)

	return func(c *fiber.Ctx) error {
		ctx, span := tracer.Start(c.Context(), fmt.Sprintf("%s %s", c.Method(), c.Path()), trace.WithAttributes(label.String("HTTP Method", c.Method())))
		c.Locals("ctx", ctx)
		defer span.End()
		return c.Next()

	}
}
