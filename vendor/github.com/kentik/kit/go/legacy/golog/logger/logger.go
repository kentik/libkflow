// logger.go: extends the log package so that each log message is
// prefixed with a string
//
// Copyright (c) 2011-2012 CloudFlare, Inc.

package logger

import (
	"bytes"
	"github.com/kentik/kit/go/legacy/common/logging"
	"runtime"
	"strings"
	"sync/atomic"
)

type Level int

var (
	// OffLogger is a dummy no-op logger.
	OffLogger = New(Levels.Off)

	// Levels is a singleton that represents possible log levels.
	Levels = struct {
		Off    Level
		Panic  Level
		Error  Level
		Warn   Level
		Info   Level
		Debug  Level
		Access Level
	}{
		Access: -1,
		Off:    0,
		Panic:  1,
		Error:  2,
		Warn:   3,
		Info:   4,
		Debug:  5,
	}

	// levelMap maps Level objects to the pretty printed name
	levelMap = map[Level]string{
		Levels.Access: "Access",
		Levels.Off:    "Off",
		Levels.Panic:  "Panic",
		Levels.Error:  "Error",
		Levels.Warn:   "Warn",
		Levels.Info:   "Info",
		Levels.Debug:  "Debug",
	}

	// CfgLevels maps strings to Level. The intent is to use this during config
	// time.
	CfgLevels = map[string]Level{
		"access": Levels.Access,
		"off":    Levels.Off,
		"panic":  Levels.Panic,
		"error":  Levels.Error,
		"warn":   Levels.Warn,
		"info":   Levels.Info,
		"debug":  Levels.Debug,
	}

	logCount  uint64 // number of messages attempted on all loggers
	dropCount uint64 // number of messages dropped on all loggers
	errCount  uint64 // number of errors seen across all loggers
)

// Stats returns the current status of the logger. It reports:
// logs: number of logs attempted to be written since startup
// pending: number of logs queued to be written
// drop: numer of logs that have been dropped, because the write queue is full, since startup
// errs: number of errors seen while trying to write logs since startup
func Stats() (logs, pending, drop, errs uint64) {
	return atomic.LoadUint64(&logCount), uint64(len(messages)), atomic.LoadUint64(&dropCount), atomic.LoadUint64(&errCount)
}

type Logger struct {
	level                Level
	sample, sampleCount  uint64 // counters to allow us to sample every "sample" access logs
	callerSkipFrameCount int    // additional caller frames to skip while detecting file:line
}

func (level Level) String() string {
	return levelMap[level]
}

func New(level Level) (l *Logger) {
	l = new(Logger)
	l.level = level
	l.sample = 1
	l.callerSkipFrameCount = 0

	return
}

// WithCallerWithSkipFrameCount creates a child logger with custom number of frames to skip.
//
// This method is intended to be used by log wrappers in order retrieve the correct file:line
func (l *Logger) WithCallerWithSkipFrameCount(skipFrameCount int) *Logger {
	instance := new(Logger)
	instance.level = l.level
	instance.sample = l.sample
	instance.callerSkipFrameCount = skipFrameCount
	return instance
}

func (l *Logger) log(level Level, prefix, format string, v []interface{}, tee bool) {
	switch {
	case l == nil:
		return
	case level == Levels.Access:
		count := atomic.AddUint64(&l.sampleCount, 1)
		if l.sample == 0 || count%l.sample != 0 {
			return
		}
	case level > l.level, level == Levels.Off:
		return
	}

	_, file, line, _ := runtime.Caller(2 + l.callerSkipFrameCount)
	caller := logCaller{stripFile(file), line}
	_ = queueMsg(&logEntry{level, prefix, format, v, caller, tee})
	// TODO: instead of ignoring error from queueMsg(), send it to stderr|stdout?
}

func (l *Logger) Printf(level Level, prefix, format string, v ...interface{}) {
	l.log(level, prefix, format, v, true)
}

// Debug logs a printf-style debug message (deprecated, please use Debugf)
func (l *Logger) Debug(prefix, format string, v ...interface{}) {
	l.log(Levels.Debug, prefix, format, v, true)
}

// Debugf logs a printf-style debug message
func (l *Logger) Debugf(prefix, format string, v ...interface{}) {
	l.log(Levels.Debug, prefix, format, v, true)
}

// Info logs a printf-style info message (deprecated, please use Infof)
func (l *Logger) Info(prefix, format string, v ...interface{}) {
	l.log(Levels.Info, prefix, format, v, true)
}

// Infof logs a printf-style info message
func (l *Logger) Infof(prefix, format string, v ...interface{}) {
	l.log(Levels.Info, prefix, format, v, true)
}

// Warn logs a printf-style warn message (deprecated, please use Warnf)
func (l *Logger) Warn(prefix, format string, v ...interface{}) {
	l.log(Levels.Warn, prefix, format, v, true)
}

// Warnf logs a printf-style warn message
func (l *Logger) Warnf(prefix, format string, v ...interface{}) {
	l.log(Levels.Warn, prefix, format, v, true)
}

// Error logs a printf-style error message (deprecated, please use Errorf)
func (l *Logger) Error(prefix, format string, v ...interface{}) {
	l.log(Levels.Error, prefix, format, v, true)
}

// Errorf logs a printf-style error message
func (l *Logger) Errorf(prefix, format string, v ...interface{}) {
	l.log(Levels.Error, prefix, format, v, true)
}

// Panic logs a printf-style panic message (deprecated, please use Panicf)
func (l *Logger) Panic(prefix, format string, v ...interface{}) {
	l.log(Levels.Panic, prefix, format, v, true)
}

// Panicf logs a printf-style panic message
func (l *Logger) Panicf(prefix, format string, v ...interface{}) {
	l.log(Levels.Panic, prefix, format, v, true)
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func (l *Logger) Level() Level {
	return l.level
}

func (l *Logger) SetAccessLogSample(sample uint64) {
	atomic.StoreUint64(&l.sample, sample)
}

func (l *Logger) Write(p []byte) (int, error) {
	level := Levels.Info
	if bytes.Contains(p, []byte("Error")) {
		level = Levels.Error
	} else if bytes.Contains(p, []byte("Warn")) {
		level = Levels.Warn
	}
	v := []interface{}{string(p)}
	l.log(level, "", "%s", v, true)

	return len(p), nil
}

// Verify Interface Compliance
var _ logging.Logger = (*Logger)(nil)

func stripFile(file string) string {
	paths := []string{
		// Most to least specific
		"vendor/github.com/kentik/",
		"vendor/github.com/",
		"vendor/",
		"build/input/",
	}

	for _, s := range paths {
		if idx := strings.Index(file, s); idx >= 0 {
			file = file[idx+len(s):]
			break
		}
	}
	return file
}

func LogNoTee(level Level, prefix string, format string, v ...interface{}) {
	New(Levels.Info).log(level, prefix, format, v, false)
}

type Child interface {
	WithCallerWithSkipFrameCount(skipFrameCount int) *Logger
}
