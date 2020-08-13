// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logging

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// Level is the underlying type of the level.
type Level uint32

// These are the different logging levels. You can set the the logging level on
// your instance of logger.
const (
	// SilentLevel, there is no log output.
	SilentLevel = iota
	// ErrorLevel, used for errors that should be addressed immediately.
	ErrorLevel
	// WarnLevel, non-critical errors that should be looked at.
	WarnLevel
	// InfoLevel, general information about what is going on in the application.
	InfoLevel
	// DebugLevel, highly verbose information usually used for debugging.
	DebugLevel
	// TraceLevel, maximum verbosity, finer grained than the Debug level.
	TraceLevel
)

// Map Levels to logrus levels.
var levelMap = map[Level]logrus.Level{
	ErrorLevel: logrus.ErrorLevel,
	WarnLevel:  logrus.WarnLevel,
	InfoLevel:  logrus.InfoLevel,
	DebugLevel: logrus.DebugLevel,
	TraceLevel: logrus.TraceLevel,
}

// AllLevels lists all available log Levels.
var AllLevels = []Level{
	SilentLevel,
	ErrorLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
	TraceLevel,
}

// ParseLevel will parse a string representation of a log level and return a
// logging Level.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "silent":
		return SilentLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	case "trace":
		return TraceLevel, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid Level: %q", lvl)
}

// UnmarshalText implements the Unmarshaller interface for the Level.
func (l *Level) UnmarshalText(text []byte) error {
	lvl, err := ParseLevel(string(text))
	if err != nil {
		return err
	}

	*l = lvl

	return nil
}

// MarshalText implements the Marshaller interface for the Level.
func (l Level) MarshalText() ([]byte, error) {
	switch l {
	case SilentLevel:
		return []byte("silent"), nil
	case ErrorLevel:
		return []byte("error"), nil
	case WarnLevel:
		return []byte("warn"), nil
	case InfoLevel:
		return []byte("info"), nil
	case DebugLevel:
		return []byte("debug"), nil
	case TraceLevel:
		return []byte("trace"), nil
	}

	return nil, fmt.Errorf("not a valid log level: %d", l)
}

// String will convert the level to a string (eg. ErrorLevel -> "error").
func (l Level) String() string {
	b, err := l.MarshalText()
	if err != nil {
		return "unknown"
	}
	return string(b)
}
