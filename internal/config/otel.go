package config

import (
	"context"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func InitTracer(config Config) func(ctx context.Context) error {
	projectID := config.ProjectID

	exporter, err := texporter.New(texporter.WithProjectID(projectID))
	if err != nil {
		l.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	otel.SetTracerProvider(tp)

	return tp.ForceFlush
}
