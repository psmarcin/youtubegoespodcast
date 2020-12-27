package config

import (
	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
)

func InitTracer(config Config) func() {
	projectID := config.ProjectID

	traceProvider, flush, err := cloudtrace.InstallNewPipeline(
		[]cloudtrace.Option{cloudtrace.WithProjectID(projectID)},
	)
	if err != nil {
		l.Fatal(err)
	}

	otel.SetTracerProvider(traceProvider)

	return flush
}
