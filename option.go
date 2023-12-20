package entslog

import (
	"context"
	"log/slog"
)

type Option struct {
	DefaultLevel slog.Leveler
	ErrorLevel   slog.Leveler
	HandleError  bool
	FilterAttrs  []string
}

var DefaultOption = Option{}

func (o Option) handle(logger *slog.Logger) *Handler {
	if o.DefaultLevel == nil {
		o.DefaultLevel = slog.LevelDebug
	}

	if o.ErrorLevel == nil {
		o.ErrorLevel = slog.LevelError
	}

	filterHandle := func(attrs ...slog.Attr) []slog.Attr {
		return attrs
	}
	if len(o.FilterAttrs) > 0 {
		filterHandle = func(attrs ...slog.Attr) []slog.Attr {
			//range attrs checking the keywords,if the keyword is in the config.FilterAttrs,then remove it
			var newAttrs []slog.Attr
			for _, attr := range attrs {
				for _, keyword := range o.FilterAttrs {
					if attr.Key == keyword {
						continue
					}
				}
				newAttrs = append(newAttrs, attr)
			}
			return newAttrs
		}
	}

	return &Handler{
		logger: logger,
		log:    func(ctx context.Context, msg string, attrs ...slog.Attr) {},
		errorHandle: func(ctx context.Context, msg string, err error) error {
			return err
		},
		filterHandle: filterHandle,
		option:       o,
	}
}
