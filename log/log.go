package log

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

type Logger zerolog.Logger

var logger Logger

func init() {
	writer := zerolog.ConsoleWriter{
		NoColor:    true,
		Out:        os.Stderr,
		TimeFormat: "2006-01-02 15:04:05",
	}

	logger = Logger(zerolog.New(writer).With().Timestamp().Logger())
}

func Global() *Logger {
	return &logger
}

func SetVerbose(level int) {
	switch level {
	case 1:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case 2:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func Debugf(format string, v ...interface{}) {
	(*zerolog.Logger)(&logger).Debug().Msg(fmt.Sprintf(format, v...))
}

func Fatal(err error) {
	(*zerolog.Logger)(&logger).Fatal().Msg(err.Error())
}

func (l *Logger) Debug(msg string, kv ...interface{}) {
	format((*zerolog.Logger)(l).Debug(), kv...).Msg(msg)
}

func (l *Logger) Info(msg string, kv ...interface{}) {
	format((*zerolog.Logger)(l).Info(), kv...).Msg(msg)
}

func (l *Logger) Warn(msg string, kv ...interface{}) {
	format((*zerolog.Logger)(l).Warn(), kv...).Msg(msg)
}

func (l *Logger) Error(msg string, kv ...interface{}) {
	format((*zerolog.Logger)(l).Error(), kv...).Msg(msg)
}

func format(event *zerolog.Event, kv ...interface{}) *zerolog.Event {
	for i := 0; i < len(kv)-1; i += 2 {
		key := kv[i].(string)
		val := kv[i+1]
		switch v := val.(type) {
		case string:
			event = event.Str(key, v)
		case int:
			event = event.Int(key, v)
		case uint:
			event = event.Uint(key, v)
		case fmt.Stringer:
			event = event.Str(key, v.String())
		default:
			event = event.Str(key, fmt.Sprintf("%v", v))
		}

	}
	return event
}
