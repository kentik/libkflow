package logging

import "github.com/rs/zerolog"

const prefixKey = "prefix"

type wrapper struct {
	l zerolog.Logger
}

func (w *wrapper) Debugf(prefix, format string, v ...interface{}) {
	w.l.Debug().Str(prefixKey, prefix).Msgf(format, v...)
}

func (w *wrapper) Infof(prefix, format string, v ...interface{}) {
	w.l.Info().Str(prefixKey, prefix).Msgf(format, v...)
}

func (w *wrapper) Errorf(prefix, format string, v ...interface{}) {
	w.l.Error().Str(prefixKey, prefix).Msgf(format, v...)
}

func (w *wrapper) Warnf(prefix, format string, v ...interface{}) {
	w.l.Warn().Str(prefixKey, prefix).Msgf(format, v...)
}

func (w *wrapper) Panicf(prefix, format string, v ...interface{}) {
	w.l.Panic().Str(prefixKey, prefix).Msgf(format, v...)
}

func Wrap(l zerolog.Logger) Logger {

	return &wrapper{
		// Skip the wrapper
		l: l.With().CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + 1).Logger(),
	}
}
