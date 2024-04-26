package entslog

import (
    "context"
    "log/slog"
)

type Handler struct {
    logger *slog.Logger
    error  func(ctx context.Context, msg string, err error) error
    option Option
    attrs  []slog.Attr
}

func (eh *Handler) errorMethod(ctx context.Context, msg string, err error) error {
    if err != nil {
        attrs := eh.Filter(ctx, slog.Any("error", err))
        eh.logger.LogAttrs(ctx,
            eh.option.errorLevel.Level(), msg,
            attrs...,
        )
    }
    return err
}

func (eh *Handler) log(ctx context.Context, msg string, attrs ...slog.Attr) {
    attrs = eh.Filter(ctx, attrs...)
    eh.logger.LogAttrs(ctx,
        eh.option.defaultLevel.Level(), msg,
        attrs...)
}

func (eh *Handler) initWith(attrs ...slog.Attr) Handler {
    ehCopy := *eh
    ehCopy.attrs = attrs
    if ehCopy.option.handleError {
        ehCopy.error = ehCopy.errorMethod
    }
    return ehCopy
}

func (eh *Handler) with(attrs ...slog.Attr) Handler {
    ehCopy := *eh
    ehCopy.attrs = attrs
    return ehCopy
}

func (eh *Handler) Filter(ctx context.Context, attrs ...slog.Attr) []slog.Attr {
    attrs = eh.option.filter(ctx, attrs...)
    return append(eh.attrs, attrs...)
}
