package config

import (
	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel/api/global"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func InitTracer(config Config) func() {
	projectID := config.ProjectID
	l.Infof("conifg: %+v", config)

	// Create Google Cloud Trace exporter to be able to retrieve
	// the collected spans.
	traceProvider, flush, err := cloudtrace.InstallNewPipeline(
		[]cloudtrace.Option{cloudtrace.WithProjectID(projectID)},
		// For this example code we use sdktrace.AlwaysSample sampler to sample all traces.
		// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)
	if err != nil {
		l.Fatal(err)
	}

	global.SetTracerProvider(traceProvider)

	return flush
}
