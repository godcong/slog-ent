package entslog

import (
	"context"
	"log/slog"
)

type Handler struct {
	logger       *slog.Logger
	log          func(ctx context.Context, msg string, attrs ...slog.Attr)
	errorHandle  func(ctx context.Context, msg string, err error) error
	filterHandle func(attrs ...slog.Attr) []slog.Attr
	option       Option
}

func (eh *Handler) driver() Handler {
	eh.log = func(ctx context.Context, msg string, attrs ...slog.Attr) {
		eh.logger.WithGroup("driver").LogAttrs(ctx, eh.option.DefaultLevel.Level(), msg, eh.filterHandle(attrs...)...)
	}
	if eh.option.HandleError {
		eh.errorHandle = func(ctx context.Context, msg string, err error) error {
			if err != nil {
				eh.logger.WithGroup("driver").LogAttrs(ctx, eh.option.ErrorLevel.Level(), msg, slog.Any("error", err))
			}
			return err
		}
	}
	return *eh
}

func (eh *Handler) tx() Handler {
	eh.log = func(ctx context.Context, msg string, attrs ...slog.Attr) {
		eh.logger.WithGroup("Tx").LogAttrs(ctx, eh.option.DefaultLevel.Level(), msg, eh.filterHandle(attrs...)...)
	}
	if eh.option.HandleError {
		eh.errorHandle = func(ctx context.Context, msg string, err error) error {
			if err != nil {
				eh.logger.WithGroup("Tx").LogAttrs(ctx, eh.option.ErrorLevel.Level(), msg, slog.Any("error", err))
			}
			return err
		}
	}
	return *eh
}
