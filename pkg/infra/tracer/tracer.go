package tracer

import (
	"fmt"
	"github.com/caarlos0/env"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
)

const TRACER_NAME = "team10_auth_service"

var Tracer = otel.Tracer(TRACER_NAME)

type Config struct {
	ServiceName   string `env:"SERVICE_NAME" envDefault:"team10_auth"`
	JaegerAddress string `env:"JAEGER_ADDRESS" envDefault:"http://jaeger-instance-collector.observability:14268/api/traces"`
}

func InitOtel() error {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("parse metrics configuration failed: %w", err)
	}

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint(cfg.JaegerAddress)),
	)
	if err != nil {
		return err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
		)))

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return nil
}
