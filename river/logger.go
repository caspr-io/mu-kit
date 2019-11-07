package river

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/rs/zerolog"
)

type ZeroLogger struct {
	traceEnabled bool
	logger       *zerolog.Logger
}

func (l *ZeroLogger) Error(msg string, err error, fields watermill.LogFields) {
	withFields(l.logger, fields).Error().Err(err).Msg(msg)
}

func (l *ZeroLogger) Info(msg string, fields watermill.LogFields) {
	withFields(l.logger, fields).Info().Msg(msg)
}

func (l *ZeroLogger) Debug(msg string, fields watermill.LogFields) {
	withFields(l.logger, fields).Debug().Msg(msg)
}

func (l *ZeroLogger) Trace(msg string, fields watermill.LogFields) {
	if l.traceEnabled {
		withFields(l.logger, fields).Debug().Str("trace", "true").Msg(msg)
	}
}

func (l *ZeroLogger) With(fields watermill.LogFields) watermill.LoggerAdapter {
	if fields == nil || len(fields) == 0 {
		return l
	}

	newLogger := withFields(l.logger, fields)

	return &ZeroLogger{logger: newLogger}
}

func withFields(logger *zerolog.Logger, fields watermill.LogFields) *zerolog.Logger {
	if fields == nil || len(fields) == 0 {
		return logger
	}

	loggerBuilder := logger.With()
	for k, v := range fields {
		loggerBuilder = loggerBuilder.Interface(k, v)
	}

	newLogger := loggerBuilder.Logger()

	return &newLogger
}

func NewZerologLogger(logger *zerolog.Logger) watermill.LoggerAdapter {
	return &ZeroLogger{true, logger}
}
