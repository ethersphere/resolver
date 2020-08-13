// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logging

import (
	"io"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

var silencedOutputWriter = ioutil.Discard

// Logging provides functions for writing to a log different levels.
type Logging interface {
	Tracef(format string, args ...interface{})
	Trace(args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Warningf(format string, args ...interface{})
	Warning(args ...interface{})
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
	WithField(key string, value interface{}) *logrus.Entry
	WithFields(fields logrus.Fields) *logrus.Entry
	WriterLevel(logrus.Level) *io.PipeWriter
	NewEntry() *logrus.Entry
}

// Logger implements the Logging interface. It wraps a logrus Logger.
type Logger struct {
	Level Level
	impl  *logrus.Logger
}

// TextFormatOptions control formatting of logger output (as text). The options
// are passed to the logrus TextFormatter.
type TextFormatOptions struct {
	FullTimestamp bool
}

// Option applies an option to Logger.
type Option func(*Logger)

// New creates a new Logger with the appropriate output and default log level.
func New(opts ...Option) *Logger {
	logger := Logger{
		impl: logrus.New(),
	}

	// Set defaults.
	logger.SetLevel(InfoLevel)
	logger.impl.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	// Apply all options to the Logger.
	for _, o := range opts {
		o(&logger)
	}

	return &logger
}

// WithLevel applies SetLevel on the Logger.
func WithLevel(lvl Level) Option {
	return func(l *Logger) {
		l.SetLevel(lvl)
	}
}

// WithOutput applies SetWriter on the Logger.
func WithOutput(w io.Writer) Option {
	return func(l *Logger) {
		l.SetOutput(w)
	}
}

// WithTextFormat applies SetTextFormat on the Logger.
func WithTextFormat(tfo TextFormatOptions) Option {
	return func(l *Logger) {
		l.SetTextFormat(tfo)
	}
}

// SetOutput sets the log output writer.
func (l *Logger) SetOutput(w io.Writer) {
	l.impl.SetOutput(w)
}

// SetLevel sets the loging level.
func (l *Logger) SetLevel(level Level) {
	l.impl.SetLevel(levelMap[level])
	if level == SilentLevel {
		l.SetOutput(silencedOutputWriter)
	}
	l.Level = level
}

// SetTextFormat overrides the logger fomatting into text and sets the text
// formatting options.
func (l *Logger) SetTextFormat(o TextFormatOptions) {
	l.impl.Formatter = &logrus.TextFormatter{
		FullTimestamp: o.FullTimestamp,
	}
}
