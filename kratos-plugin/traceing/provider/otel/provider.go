package otel

import (
	"github.com/sirupsen/logrus"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type Config struct {
	Endpoint    string  `json:"endpoint" toml:"endpoint"`
	Fraction    float64 `json:"fraction" toml:"fraction"`
	ServiceName string  `json:"service_name" toml:"service_name"`
	Version     string  `json:"version" toml:"version"`
}

func NewTracerProvider(cfg Config) (*tracesdk.TracerProvider, error) {
	logrus.Infof("cfg %#v", cfg)
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Endpoint)))
	if err != nil {
		panic(err)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.TraceIDRatioBased(cfg.Fraction)),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),

		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.Version),
		)),
	)
	return tp, nil
}
