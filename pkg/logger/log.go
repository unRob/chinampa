// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package logger

import (
	"context"

	"git.rob.mx/nidito/chinampa/pkg/runtime"
	"github.com/sirupsen/logrus"
)

var componentKey = "_component"

func init() {
	logrus.SetFormatter(new(Formatter))
}

var Main = logrus.WithContext(context.Background())

func Sub(name string) *logrus.Entry {
	return logrus.WithField(componentKey, name)
}

type Level int

const (
	LevelPanic Level = iota
	LevelFatal
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
	LevelTrace
)

func Configure(name string, level Level) {
	Main = logrus.WithField(componentKey, name)
	if runtime.SilenceEnabled() {
		logrus.SetLevel(logrus.ErrorLevel)
	} else {
		logrus.SetLevel(logrus.AllLevels[level])
	}
}

func Debug(args ...any) {
	Main.Debug(args...)
}

func Debugf(format string, args ...any) {
	Main.Debugf(format, args...)
}

func Info(args ...any) {
	Main.Info(args...)
}

func Infof(format string, args ...any) {
	Main.Infof(format, args...)
}

func Warn(args ...any) {
	Main.Warn(args...)
}

func Warnf(format string, args ...any) {
	Main.Warnf(format, args...)
}

func Error(args ...any) {
	Main.Error(args...)
}

func Errorf(format string, args ...any) {
	Main.Errorf(format, args...)
}

func Fatal(args ...any) {
	Main.Fatal(args...)
}

func Fatalf(format string, args ...any) {
	Main.Fatalf(format, args...)
}

func Trace(args ...any) {
	Main.Trace(args...)
}

func Tracef(format string, args ...any) {
	Main.Tracef(format, args...)
}
