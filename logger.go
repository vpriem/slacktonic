// Package slacktonic provides middleware for handling Slack slash commands in Gin applications.
package slacktonic

import (
	"context"
)

// Logger is an interface for logging errors that occur during request processing.
// It is designed to be compatible with the standard library's log/slog package.
type Logger interface {
	// InfoContext logs an info message with the given context and key-value pairs.
	InfoContext(ctx context.Context, msg string, args ...any)

	// WarnContext logs a warning message with the given context and key-value pairs.
	WarnContext(ctx context.Context, msg string, args ...any)

	// ErrorContext logs an error message with the given context and key-value pairs.
	ErrorContext(ctx context.Context, msg string, args ...any)
}

// noopLogger is a Logger implementation that does nothing.
type noopLogger struct{}

func (l *noopLogger) InfoContext(_ context.Context, _ string, _ ...any) {}

func (l *noopLogger) WarnContext(_ context.Context, _ string, _ ...any) {}

func (l *noopLogger) ErrorContext(_ context.Context, _ string, _ ...any) {}
