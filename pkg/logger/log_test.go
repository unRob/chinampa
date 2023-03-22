// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package logger_test

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	. "git.rob.mx/nidito/chinampa/pkg/logger"
	rt "git.rob.mx/nidito/chinampa/pkg/runtime"
	"github.com/sirupsen/logrus"
)

func withEnv(t *testing.T, env map[string]string) {
	prevEnv := os.Environ()
	for _, entry := range prevEnv {
		parts := strings.SplitN(entry, "=", 2)
		os.Unsetenv(parts[0])
	}

	for k, v := range env {
		os.Setenv(k, v)
	}

	t.Cleanup(func() {
		rt.ResetParsedFlagsCache()

		for k := range env {
			os.Unsetenv(k)
		}
		for _, entry := range prevEnv {
			parts := strings.SplitN(entry, "=", 2)
			os.Setenv(parts[0], parts[1])
		}
	})
}

func TestFormatter(t *testing.T) {
	now := strings.Replace(time.Now().Local().Format(time.DateTime), " ", "T", 1)
	cases := []struct {
		Color   bool
		Verbose bool
		Call    func(args ...any)
		Expects string
		Level   logrus.Level
	}{
		{
			Color:   true,
			Call:    Info,
			Expects: "message",
			Level:   logrus.InfoLevel,
		},
		{
			Color:   true,
			Verbose: true,
			Call:    Info,
			Expects: fmt.Sprintf("\033[2m%s\033[0m \033[2minfo\033[0m\033[2m\033[0m	message", now),
			Level:   logrus.InfoLevel,
		},
		{
			Color:   true,
			Call:    Debug,
			Expects: "",
			Level:   logrus.InfoLevel,
		},
		{
			Call:    Debug,
			Expects: "DEBUG: message",
			Level:   logrus.DebugLevel,
		},
		{
			Color:   true,
			Call:    Debug,
			Expects: "\033[2mDEBUG: \033[0m\033[2mmessage\033[0m",
			Level:   logrus.DebugLevel,
		},
		{
			Color:   true,
			Verbose: true,
			Call:    Debug,
			Expects: fmt.Sprintf("\033[2m%s\033[0m \033[2mdebug\033[0m\033[2m\033[0m\t\033[2mmessage\033[0m",
				now),
			Level: logrus.DebugLevel,
		},
		{
			Call:    Trace,
			Expects: "",
			Level:   logrus.DebugLevel,
		},
		{
			Call:    Trace,
			Expects: "TRACE: message",
			Level:   logrus.TraceLevel,
		},
		{
			Call:    Warn,
			Expects: "",
			Level:   logrus.ErrorLevel,
		},
		{
			Call:    Warn,
			Expects: "WARNING: message",
			Level:   logrus.InfoLevel,
		},
		{
			Call:    Warn,
			Level:   logrus.InfoLevel,
			Color:   true,
			Verbose: true,
			Expects: fmt.Sprintf("\033[2m%s\033[0m \033[1;93mwarning\033[0m\033[2m\033[0m\tmessage", now),
		},
		{
			Call:    Warn,
			Level:   logrus.InfoLevel,
			Color:   true,
			Expects: "\033[1;43;30m WARNING \033[0m message",
		},
		{
			Call:    Error,
			Expects: "ERROR: message",
			Level:   logrus.ErrorLevel,
		},
		{
			Call:    Error,
			Expects: "ERROR: message",
			Level:   logrus.InfoLevel,
		},
		{
			Call:    Error,
			Level:   logrus.InfoLevel,
			Color:   true,
			Verbose: true,
			Expects: fmt.Sprintf("\033[2m%s\033[0m \033[1;91merror\033[0m\033[2m\033[0m\tmessage", now),
		},
		{
			Call:    Error,
			Level:   logrus.InfoLevel,
			Color:   true,
			Expects: "\033[1;41;225m ERROR \033[0m message",
		},
	}

	for _, c := range cases {
		fname := runtime.FuncForPC(reflect.ValueOf(c.Call).Pointer()).Name()
		comps := []string{fname, c.Level.String()}
		if c.Color {
			comps = append(comps, "color")
		}
		if c.Verbose {
			comps = append(comps, "verbose")
		}
		name := strings.Join(comps, "/")
		t.Run(name, func(t *testing.T) {
			env := map[string]string{
				"COLOR":   "",
				"VERBOSE": "",
			}
			if c.Color {
				env["COLOR"] = "always"
			} else {
				env["NO_COLOR"] = "1"
			}
			if c.Verbose {
				env["VERBOSE"] = "1"
			}
			withEnv(t, env)
			data := bytes.Buffer{}
			logrus.SetLevel(c.Level)
			logrus.SetOutput(&data)
			c.Call("message")
			expected := c.Expects
			if c.Expects != "" {
				expected = c.Expects + "\n"
			}

			if res := data.String(); res != expected {
				t.Fatalf("%s:\ngot   : %s\nwanted: %v", name, res, expected)
			}
		})
	}
}
