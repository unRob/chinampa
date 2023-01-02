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

func DebugEnabled() bool {
	return isTrueIsh(os.Getenv(env.Debug))
}

func ValidationEnabled() bool {
	return isFalseIsh(os.Getenv(env.ValidationDisabled))
}

func VerboseEnabled() bool {
	return isTrueIsh(os.Getenv(env.Verbose))
}

func SilenceEnabled() bool {
	for _, arg := range os.Args {
		if arg == "--silent" {
			return true
		}
	}

	if VerboseEnabled() {
		return false
	}

	return isTrueIsh(os.Getenv(env.Silent))
}

func ColorEnabled() bool {
	return isFalseIsh(os.Getenv(env.NoColor)) && !UnstyledHelpEnabled()
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
