package ports

import (
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/semconv"
)

func RequestContext() func(c *fiber.Ctx) error {
	tracer := global.TracerProvider().Tracer(
		"yt.psmarcin.dev/api",
		trace.WithInstrumentationVersion("1.0.0"),
	)

	return func(c *fiber.Ctx) error {
		ctx, span := tracer.Start(
			c.Context(),
			"yt.psmarcin.dev/http/request",
			trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(semconv.HTTPMethodKey.String(c.Method())),
			trace.WithAttributes(semconv.HTTPUrlKey.String(c.OriginalURL())),
		)
		c.Locals("ctx", ctx)
		defer span.End()
		c.Next()

		span.SetAttributes(semconv.HTTPAttributesFromHTTPStatusCode(c.Response().StatusCode())...)

		return nil
	}
}
