// Copyright (c) 2024 OrigAdmin. All rights reserved.

// Package entslog for entgo.io/ent
package entslog

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type (
	// TraceFunc generates a unique identifier for a log entry to string.
	TraceFunc func(context.Context) string
	// FilterAttrs defines a function to filter out attributes from log entries.
	FilterAttrs func(context.Context, ...slog.Attr) []slog.Attr
	// Option defines configuration options for the logging handler.
	Option struct {
		handleError bool         // HandleError determines whether errors encountered during logging are handled.
		logger      *slog.Logger // Logger specifies the logger to be used for logging.
		level       slog.Leveler // DefaultLevel specifies the default log level for messages.
		errorLevel  slog.Leveler // ErrorLevel specifies the log level for error messages.
		trace       TraceFunc    // GenerateID is a function to generate unique IDs for log entries.
		filter      FilterAttrs  // Filters specifies the set of attributes to filter out from logged messages.
	}
	// Setting is a type alias for the settings.Setting type.
	Setting = func(*Option)
)

// defaultOption provides the default configuration options for the logging handler.
var defaultOption = Option{
	logger:      slog.Default(),  // Defaults to the default logger.
	level:       slog.LevelInfo,  // Defaults to Info level.
	errorLevel:  slog.LevelError, // Defaults to Error level.
	handleError: true,            // Defaults to handling errors.
	filter:      emptyFilter,     // Defaults to no filtering.
	trace:       traceUUID,       // Uses the package-level trace function to generate log entry IDs by default.
}

func emptyFilter(_ context.Context, attrs ...slog.Attr) []slog.Attr {
	return attrs
}

// traceUUID generates a unique identifier for a log entry using UUIDs.
func traceUUID(ctx context.Context) string {
	return uuid.Must(uuid.NewRandom()).String()
}

// WithDefaultLevel sets the default log level for the given logging options.
//
// - `level`: The log level to be set as the default.
//
// Returns a function that accepts an `*Option` parameter, modifies it by setting the default log level,
// and returns the updated `*Option` pointer.
func WithDefaultLevel(level slog.Leveler) Setting {
	return func(o *Option) {
		o.level = level
	}
}

// WithErrorLevel sets the error log level and enables handling of errors for the given logging options.
//
// - `level`: The log level to be set for error logging.
//
// Returns a function that accepts an `*Option` parameter, modifies it by setting the error log level
// and enabling error handling, then returns the updated `*Option` pointer.
func WithErrorLevel(level slog.Leveler) Setting {
	return func(option *Option) {
		option.errorLevel = level
		option.handleError = true
	}
}

// WithError explicitly enables or disables error handling for the given logging options.
//
// - `handleError`: A boolean indicating whether to enable (true) or disable (false) error handling.
//
// Returns a function that accepts an `*Option` parameter, modifies it by setting the error handling flag,
// and returns the updated `*Option` pointer.
func WithError() Setting {
	return func(option *Option) {
		option.handleError = true
	}
}

// WithFilter specifies a list of attribute names to filter out from logged messages for the given logging options.
//
// - `attrs`: A variadic list of strings representing attribute names to be filtered.
//
// Returns a function that accepts an `*Option` parameter, modifies it by setting the list of filtered attributes,
// and returns the updated `*Option` pointer.
func WithFilter(filter func(context.Context, ...slog.Attr) []slog.Attr) Setting {
	return func(option *Option) {
		option.filter = filter
	}
}

// WithTrace assigns a custom ID generation function for the given logging options.
// This function will be used to generate unique IDs for log entries within a given context.
//
// - `trace`: A function that accepts a `context.Context` and returns a string representing the generated ID.
//
// Returns a function that accepts an `*Option` parameter, modifies it by setting the custom ID generation function,
// and returns the updated `*Option` pointer.
func WithTrace(trace func(context.Context) string) Setting {
	return func(option *Option) {
		option.trace = trace
	}
}

// WithLogger specifies the logger to be used for logging.
// If not specified, the default logger will be used.
//
// - `logger`: The logger to be used for logging.
//
// Returns a function that accepts an `*Option` parameter, modifies it by setting the logger,
// and returns the updated `*Option` pointer.
func WithLogger(logger *slog.Logger) Setting {
	return func(option *Option) {
		option.logger = logger
	}
}

// make configures and returns a new logging handler based on the provided options.
