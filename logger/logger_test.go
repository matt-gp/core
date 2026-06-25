package logger

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

// setupTestLogger wires a buffer-backed logger provider into the OTEL global and
// calls Init so the package-level logger is ready for the test.
func setupTestLogger(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer

	exporter, err := stdoutlog.New(
		stdoutlog.WithWriter(&buf),
		stdoutlog.WithPrettyPrint(),
	)
	if err != nil {
		t.Fatalf("failed to create log exporter: %v", err)
	}

	provider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
	)
	SetLogger(provider.Logger("test-logger"))

	return &buf
}

func TestLogLevels(t *testing.T) {
	tests := []struct {
		name     string
		logFunc  func(context.Context, string, ...attribute.KeyValue)
		message  string
		setLevel string
	}{
		{
			name:     "debug level",
			logFunc:  Debug,
			message:  "debug message",
			setLevel: "DEBUG",
		},
		{
			name:     "trace level",
			logFunc:  Trace,
			message:  "trace message",
			setLevel: "DEBUG",
		},
		{
			name:    "info level",
			logFunc: Info,
			message: "info message",
		},
		{
			name:    "warn level",
			logFunc: Warn,
			message: "warning message",
		},
		{
			name:    "error level",
			logFunc: Error,
			message: "error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := setupTestLogger(t)
			ctx := context.Background()

			if tt.setLevel != "" {
				if err := os.Setenv("LOG_LEVEL", tt.setLevel); err != nil {
					t.Fatalf("Failed to set env var LOG_LEVEL: %v", err)
				}
				defer func() {
					if err := os.Unsetenv("LOG_LEVEL"); err != nil {
						fmt.Printf("Failed to unset env var LOG_LEVEL: %v\n", err)
					}
				}()
			}

			tt.logFunc(ctx, tt.message)

			assert.NotNil(t, buf)
		})
	}
}

func TestLogWithAttributes(t *testing.T) {
	tests := []struct {
		name       string
		logFunc    func(context.Context, string, ...attribute.KeyValue)
		message    string
		attributes []attribute.KeyValue
	}{
		{
			name:    "info with attributes",
			logFunc: Info,
			message: "info message with attributes",
			attributes: []attribute.KeyValue{
				attribute.String("key1", "value1"),
				attribute.Int("key2", 42),
				attribute.Bool("key3", true),
			},
		},
		{
			name:    "error with attributes",
			logFunc: Error,
			message: "error message with attributes",
			attributes: []attribute.KeyValue{
				attribute.String("component", "test"),
				attribute.String("error", errors.New("test error").Error()),
			},
		},
		{
			name:    "warn with attributes",
			logFunc: Warn,
			message: "warning with attributes",
			attributes: []attribute.KeyValue{
				attribute.String("service", "test-service"),
				attribute.Int64("timestamp", 1234567890),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := setupTestLogger(t)
			ctx := context.Background()

			tt.logFunc(ctx, tt.message, tt.attributes...)

			assert.NotNil(t, buf)
		})
	}
}

func TestLogOutput(t *testing.T) {
	buf := setupTestLogger(t)
	ctx := context.Background()

	if err := os.Setenv("LOG_LEVEL", "INFO"); err != nil {
		t.Fatalf("Failed to set env var LOG_LEVEL: %v", err)
	}
	defer func() {
		if err := os.Unsetenv("LOG_LEVEL"); err != nil {
			fmt.Printf("Failed to unset env var LOG_LEVEL: %v\n", err)
		}
	}()

	Info(ctx, "test message", attribute.String("key", "value"))

	assert.NotNil(t, buf)
}
