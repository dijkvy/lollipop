package otel

import (
	"context"
	"testing"
)

func TestNewTracerProvider(t *testing.T) {
	provider, err := NewTracerProvider(Config{
		Version:     "demo-v1.0.0.0rc",
		Endpoint:    "http://127.0.0.1:14268/api/traces",
		ServiceName: "demo.xx.local",
		Fraction:    1,
	})
	if err != nil {
		panic(err)
	}

	defer provider.Shutdown(context.TODO())

	tracer := provider.Tracer("TestNewTracerProvider")
	ctx, span := tracer.Start(context.TODO(), "TestNewTracerProvider_inner")
	defer span.End()
	_ = ctx
}
