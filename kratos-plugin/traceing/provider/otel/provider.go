package otel

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
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
	exp, err := jaeger.NewRawExporter(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Endpoint)))
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.TraceIDRatioBased(cfg.Fraction)),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.Version),
		)),
	)
	return tp, nil
}
