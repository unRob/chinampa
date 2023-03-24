// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0

/*
Package runtime presents environment related information useful during your program's runtime.
*/
package runtime

import (
	"os"
	"strconv"
	"strings"

	"git.rob.mx/nidito/chinampa/pkg/env"
)

// Executable is the name of our binary, and should be set
// using `chinampa.Execute(config chinampa.Config)`.
var Executable = "chinampa"

var falseIshValues = []string{
	"",
	"0",
	"no",
	"false",
	"disable",
	"disabled",
	"off",
	"never",
}

var trueIshValues = []string{
	"1",
	"yes",
	"true",
	"enable",
	"enabled",
	"on",
	"always",
}

func isFalseIsh(val string) bool {
	for _, negative := range falseIshValues {
		if val == negative {
			return true
		}
	}

	return false
}

func isTrueIsh(val string) bool {
	for _, positive := range trueIshValues {
		if val == positive {
			return true
		}
	}

	return false
}

var _flags map[string]bool

// ResetParsedFlagsCache resets the cached parsed global flags.
func ResetParsedFlagsCache() {
	_flags = nil
}

func flagInArgs(name string) bool {
	if _flags == nil {
		_flags = map[string]bool{}
		for _, arg := range os.Args {
			switch arg {
			case "--verbose":
				_flags["verbose"] = true
				delete(_flags, "silent")
			case "--silent":
				_flags["silent"] = true
				delete(_flags, "verbose")
			case "--color":
				_flags["color"] = true
				delete(_flags, "no-color")
			case "--no-color":
				_flags["no-color"] = true
				delete(_flags, "color")
			case "--skip-validation":
				_flags["skip-validation"] = true
			}
		}
	}

	_, ok := _flags[name]
	return ok
}

// DebugEnabled tells if debugging was requested.
func DebugEnabled() bool {
	return isTrueIsh(os.Getenv(env.Debug))
}

// DebugEnabled tells if debugging was requested.
func ValidationEnabled() bool {
	return !flagInArgs("skip-validation") && isFalseIsh(os.Getenv(env.ValidationDisabled))
}

// VerboseEnabled tells if verbose output was requested.
func VerboseEnabled() bool {
	if flagInArgs("silent") {
		return false
	}
	return isTrueIsh(os.Getenv(env.Verbose)) || flagInArgs("verbose")
}

// SilenceEnabled tells if silencing of output was requested.
func SilenceEnabled() bool {
	if flagInArgs("verbose") {
		return false
	}
	if flagInArgs("silent") {
		return true
	}

	return isTrueIsh(os.Getenv(env.Silent)) || flagInArgs("silent")
}

// ColorEnabled tells if colorizing output was requested.
func ColorEnabled() bool {
	if flagInArgs("color") {
		return true
	}

	// we're talking to ttys, we want color unless NO_COLOR/--no-color
	return !(isTrueIsh(os.Getenv(env.NoColor)) || flagInArgs("no-color"))
}

// HelpStyle returns the style to use when rendering help.
func HelpStyle() string {
	return strings.ToLower(os.Getenv(env.HelpStyle))
}

// EnvironmentMap returns a map of environment keys for color, debugging and verbosity and their values, ready for `os.Setenv`.
func EnvironmentMap() map[string]string {
	res := map[string]string{}
	trueString := strconv.FormatBool(true)

	if !ColorEnabled() {
		res[env.NoColor] = trueString
	} else if isTrueIsh(os.Getenv(env.ForceColor)) {
		res[env.ForceColor] = "always"
	}

	if DebugEnabled() {
		res[env.Debug] = trueString
	}

	if VerboseEnabled() {
		res[env.Verbose] = trueString
	} else if SilenceEnabled() {
		res[env.Silent] = trueString
	}

	return res
}
