// Package logger provides convenience functions for OpenTelemetry logging.
package logger

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
)

var loggerProvider log.Logger

func init() {
	// Initialize with a default logger from the global provider
	// This will be a no-op logger unless a global provider has been set
	loggerProvider = global.GetLoggerProvider().Logger("default")
}

// SetProvider configures the OTEL logger provider.
// Call this once after otel.NewProvider, passing the instrumentation scope name for your service.
func SetProvider(provider log.Logger) {
	loggerProvider = provider
}

// Debug emits a debug log using OpenTelemetry logging.
func Debug(ctx context.Context, msg string, attrs ...attribute.KeyValue) {
	record := log.Record{}
	record.SetSeverity(log.SeverityDebug)
	record.SetBody(log.StringValue(msg))
	for _, attr := range attrs {
		record.AddAttributes(log.KeyValueFromAttribute(attr))
	}
	loggerProvider.Emit(ctx, record)
}

// Trace emits a trace log using OpenTelemetry logging.
func Trace(ctx context.Context, msg string, attrs ...attribute.KeyValue) {
	record := log.Record{}
	record.SetSeverity(log.SeverityTrace)
	record.SetBody(log.StringValue(msg))
	for _, attr := range attrs {
		record.AddAttributes(log.KeyValueFromAttribute(attr))
	}
	loggerProvider.Emit(ctx, record)
}

// Info emits an info log using OpenTelemetry logging.
func Info(ctx context.Context, msg string, attrs ...attribute.KeyValue) {
	record := log.Record{}
	record.SetSeverity(log.SeverityInfo)
	record.SetBody(log.StringValue(msg))
	for _, attr := range attrs {
		record.AddAttributes(log.KeyValueFromAttribute(attr))
	}
	loggerProvider.Emit(ctx, record)
}

// Warn emits a warning log using OpenTelemetry logging.
func Warn(ctx context.Context, msg string, attrs ...attribute.KeyValue) {
	record := log.Record{}
	record.SetSeverity(log.SeverityWarn)
	record.SetBody(log.StringValue(msg))
	for _, attr := range attrs {
		record.AddAttributes(log.KeyValueFromAttribute(attr))
	}
	loggerProvider.Emit(ctx, record)
}

// Error emits an error log using OpenTelemetry logging.
func Error(ctx context.Context, msg string, attrs ...attribute.KeyValue) {
	record := log.Record{}
	record.SetSeverity(log.SeverityError)
	record.SetBody(log.StringValue(msg))
	for _, attr := range attrs {
		record.AddAttributes(log.KeyValueFromAttribute(attr))
	}
	loggerProvider.Emit(ctx, record)
}
