// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package runtime

import (
	"os"
	"strconv"

	"git.rob.mx/nidito/chinampa/pkg/env"
)

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

func ResetParsedFlags() {
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

func DebugEnabled() bool {
	return isTrueIsh(os.Getenv(env.Debug))
}

func ValidationEnabled() bool {
	return !flagInArgs("skip-validation") && isFalseIsh(os.Getenv(env.ValidationDisabled))
}

func VerboseEnabled() bool {
	if flagInArgs("silent") {
		return false
	}
	return isTrueIsh(os.Getenv(env.Verbose)) || flagInArgs("verbose")
}

func SilenceEnabled() bool {
	if flagInArgs("verbose") {
		return false
	}
	if flagInArgs("silent") {
		return true
	}

	return isTrueIsh(os.Getenv(env.Silent)) || flagInArgs("silent")
}

func ColorEnabled() bool {
	if flagInArgs("color") {
		return true
	}

	// we're talking to ttys, we want color unless NO_COLOR/--no-color
	return !(isTrueIsh(os.Getenv(env.NoColor)) || UnstyledHelpEnabled() || flagInArgs("no-color"))
}

func UnstyledHelpEnabled() bool {
	return isTrueIsh(os.Getenv(env.HelpUnstyled))
}

// EnvironmentMap returns the resolved environment map.
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
