package otelwrapper

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"detectviz/pkg/shared/log"
)

// initializeResource initializes the OpenTelemetry resource.
// zh: initializeResource 初始化 OpenTelemetry 資源。
func (w *OtelWrapper) initializeResource(ctx context.Context) error {
	attrs := []attribute.KeyValue{
		semconv.ServiceName(w.config.ServiceName),
		semconv.ServiceVersion(w.config.ServiceVersion),
		semconv.DeploymentEnvironment(w.config.Environment),
	}

	// Add custom attributes
	if w.config.Resource.Organization != "" {
		attrs = append(attrs, attribute.String("organization", w.config.Resource.Organization))
	}
	if w.config.Resource.Team != "" {
		attrs = append(attrs, attribute.String("team", w.config.Resource.Team))
	}

	// Add custom attributes from config
	for k, v := range w.config.Resource.CustomAttrs {
		attrs = append(attrs, attribute.String(k, v))
	}

	// Add global attributes
	for k, v := range w.config.Attributes {
		attrs = append(attrs, attribute.String(k, v))
	}

	var resourceOptions []resource.Option
	resourceOptions = append(resourceOptions, resource.WithAttributes(attrs...))

	// Detect environment if enabled
	if w.config.Resource.DetectHost {
		resourceOptions = append(resourceOptions, resource.WithHost())
	}
	if w.config.Resource.DetectProcess {
		resourceOptions = append(resourceOptions, resource.WithProcess())
	}
	if w.config.Resource.DetectRuntime {
		resourceOptions = append(resourceOptions, resource.WithProcessRuntimeName(),
			resource.WithProcessRuntimeVersion(), resource.WithProcessRuntimeDescription())
	}

	res, err := resource.New(ctx, resourceOptions...)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	w.resource = res

	log.L(ctx).Info("OpenTelemetry resource initialized",
		"service_name", w.config.ServiceName,
		"service_version", w.config.ServiceVersion,
		"environment", w.config.Environment)

	return nil
}

// initializeTracing initializes the tracing provider.
// zh: initializeTracing 初始化追蹤提供器。
func (w *OtelWrapper) initializeTracing(ctx context.Context) error {
	// For now, create a basic tracer provider
	// Full implementation would include exporters based on config
	w.tracerProvider = otel.GetTracerProvider()

	log.L(ctx).Info("OpenTelemetry tracing initialized")
	return nil
}

// initializeMetrics initializes the metrics provider.
// zh: initializeMetrics 初始化指標提供器。
func (w *OtelWrapper) initializeMetrics(ctx context.Context) error {
	// For now, create a basic meter provider
	// Full implementation would include exporters based on config
	w.meterProvider = otel.GetMeterProvider()

	log.L(ctx).Info("OpenTelemetry metrics initialized")
	return nil
}

// initializeLogging initializes logging integration.
// zh: initializeLogging 初始化日誌整合。
func (w *OtelWrapper) initializeLogging(ctx context.Context) error {
	// Logger provider will be injected from the logger plugin
	// This method can be used for additional logging setup if needed

	log.L(ctx).Info("OpenTelemetry logging integration initialized")
	return nil
}

// setupPropagators sets up the global propagators.
// zh: setupPropagators 設定全域傳播器。
func (w *OtelWrapper) setupPropagators() {
	// Set up composite propagator with TraceContext and Baggage
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
}
