package tntlogger

import (
	"context"

	"github.com/rs/zerolog"
)

type defaultLogger struct {
	logger zerolog.Logger
}

func NewTntLogger(lg *zerolog.Logger) *defaultLogger {
	return &defaultLogger{
		logger: lg.With().Str("source", "tarantool").Logger(),
	}
}

func (d defaultLogger) Debugf(_ context.Context, format string, v ...any) {
	d.logger.Debug().Msgf(format, v...)
}

func (d defaultLogger) Infof(_ context.Context, format string, v ...any) {
	d.logger.Info().Msgf(format, v...)
}

func (d defaultLogger) Warnf(_ context.Context, format string, v ...any) {
	d.logger.Warn().Msgf(format, v...)
}

func (d defaultLogger) Errorf(_ context.Context, format string, v ...any) {
	d.logger.Error().Msgf(format, v...)
}
