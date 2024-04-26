// Package entslog provides logging capabilities initWith ent entities.
package entslog

import (
    "context"
    "log/slog"

    "github.com/google/uuid"
)

// Option defines configuration options for the logging handler.
type Option struct {
    defaultLevel slog.Leveler // DefaultLevel specifies the default log level for messages.
    errorLevel   slog.Leveler // ErrorLevel specifies the log level for error messages.
    handleError  bool         // HandleError determines whether errors encountered during logging are handled.
    filter       func(context.Context,
    ...slog.Attr) []slog.Attr                   // Filters specifies the set of attributes to filter out from logged messages.
    generateID func(ctx context.Context) string // GenerateID is a function to generate unique IDs for log entries.
}

// DefaultOption provides the default configuration options for the logging handler.
var DefaultOption = Option{
    defaultLevel: slog.LevelInfo,  // Defaults to Info level.
    errorLevel:   slog.LevelError, // Defaults to Error level.
    handleError:  true,            // Defaults to handling errors.
    filter: func(ctx context.Context,
            attrs ...slog.Attr) []slog.Attr {
        return attrs
    },                      // Defaults to no filtering.
    generateID: generateID, // Uses the package-level generateID function to generate log entry IDs by default.
}

func settings(opts ...func(*Option) *Option) *Option {
    option := DefaultOption
    for _, opt := range opts {
        opt(&option)
    }
    return &option
}

// WithDefaultLevel sets the default log level for the given logging options.
//
// - `level`: The log level to be set as the default.
//
// Returns a function that accepts an `*Option` parameter, modifies it by setting the default log level,
// and returns the updated `*Option` pointer.
func WithDefaultLevel(level slog.Leveler) func(*Option) *Option {
    return func(option *Option) *Option {
        option.defaultLevel = level
        return option
    }
}

// WithErrorLevel sets the error log level and enables handling of errors for the given logging options.
//
// - `level`: The log level to be set for error logging.
//
// Returns a function that accepts an `*Option` parameter, modifies it by setting the error log level
// and enabling error handling, then returns the updated `*Option` pointer.
func WithErrorLevel(level slog.Leveler) func(*Option) *Option {
    return func(option *Option) *Option {
        option.errorLevel = level
        option.handleError = true
        return option
    }
}

// WithHandleError explicitly enables or disables error handling for the given logging options.
//
// - `handleError`: A boolean indicating whether to enable (true) or disable (false) error handling.
//
// Returns a function that accepts an `*Option` parameter, modifies it by setting the error handling flag,
// and returns the updated `*Option` pointer.
func WithHandleError(handleError bool) func(*Option) *Option {
    return func(option *Option) *Option {
        option.handleError = handleError
        return option
    }
}

// WithFilter specifies a list of attribute names to filter out from logged messages for the given logging options.
//
// - `attrs`: A variadic list of strings representing attribute names to be filtered.
//
// Returns a function that accepts an `*Option` parameter, modifies it by setting the list of filtered attributes,
// and returns the updated `*Option` pointer.
func WithFilter(filter func(context.Context, ...slog.Attr) []slog.Attr) func(*Option) *Option {
    return func(option *Option) *Option {
        option.filter = filter
        return option
    }
}

// WithGenerateID assigns a custom ID generation function for the given logging options.
// This function will be used to generate unique IDs for log entries within a given context.
//
// - `generateID`: A function that accepts a `context.Context` and returns a string representing the generated ID.
//
// Returns a function that accepts an `*Option` parameter, modifies it by setting the custom ID generation function,
// and returns the updated `*Option` pointer.
func WithGenerateID(generateID func(context.Context) string) func(*Option) *Option {
    return func(option *Option) *Option {
        option.generateID = generateID
        return option
    }
}

// generateID generates a unique identifier for a log entry using UUIDs.
func generateID(ctx context.Context) string {
    return uuid.New().String()
}

// handle configures and returns a new logging handler based on the provided options.
func (o Option) handle(logger *slog.Logger) *Handler {
    // Define a filter function to modify log attributes based on the FilterAttrs option.

    // Return a configured logging handler.
    return &Handler{
        logger: logger,
        error: func(ctx context.Context, msg string, err error) error {
            return err
        },
        option: o,
    }
}
