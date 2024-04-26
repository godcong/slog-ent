package entslog

import (
	"context"
	"log/slog"
)

type Handler struct {
	logger       *slog.Logger
	log          func(ctx context.Context, msg string, attrs ...slog.Attr)
	error        func(ctx context.Context, msg string, err error) error
	filters      map[string]struct{}
	filterHandle func(attrs ...slog.Attr) []slog.Attr
	option       Option
	attrs        []slog.Attr
}

func (eh *Handler) initWith(attrs ...slog.Attr) Handler {
	ehCopy := *eh
	ehCopy.attrs = attrs
	ehCopy.log = func(ctx context.Context, msg string, attrs ...slog.Attr) {
		attrs = eh.filterHandle(attrs...)
		attrs = append(ehCopy.attrs, attrs...)
		ehCopy.logger.LogAttrs(ctx,
			eh.option.defaultLevel.Level(), msg,
			attrs...)
	}
	if ehCopy.option.handleError {
		ehCopy.error = func(ctx context.Context, msg string, err error) error {
			if err != nil {
				attrs = append(ehCopy.attrs, slog.Any("error", err))
				ehCopy.logger.LogAttrs(ctx,
					eh.option.errorLevel.Level(), msg,
					attrs...,
				)
			}
			return err
		}
	}
	return ehCopy
}

func (eh *Handler) with(attrs ...slog.Attr) Handler {
	ehCopy := *eh
	ehCopy.attrs = attrs
	return ehCopy
}
