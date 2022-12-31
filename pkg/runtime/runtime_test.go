// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package runtime_test

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"testing"

	"git.rob.mx/nidito/chinampa/pkg/env"
	. "git.rob.mx/nidito/chinampa/pkg/runtime"
)

func TestEnabled(t *testing.T) {
	defer func() { os.Setenv(env.Verbose, "") }()

	cases := []struct {
		Name    string
		Func    func() bool
		Expects bool
	}{
		{
			Name:    env.Verbose,
			Func:    VerboseEnabled,
			Expects: true,
		},
		{
			Name:    env.Verbose,
			Func:    SilenceEnabled,
			Expects: false,
		},
		{
			Name: env.ValidationDisabled,
			Func: ValidationEnabled,
		},
		{
			Name: env.NoColor,
			Func: ColorEnabled,
		},
		{
			Name: env.HelpUnstyled,
			Func: ColorEnabled,
		},
		{
			Name:    env.Debug,
			Func:    DebugEnabled,
			Expects: true,
		},
		{
			Name:    env.HelpUnstyled,
			Func:    UnstyledHelpEnabled,
			Expects: true,
		},
	}

	for _, c := range cases {
		fname := runtime.FuncForPC(reflect.ValueOf(c.Func).Pointer()).Name()
		name := fmt.Sprintf("%v/%s", fname, c.Name)
		enabled := []string{
			"yes", "true", "1", "enabled",
		}
		for _, val := range enabled {
			t.Run("enabled-"+val, func(t *testing.T) {
				os.Setenv(c.Name, val)
				if c.Func() != c.Expects {
					t.Fatalf("%s wasn't enabled with a valid value: %s", name, val)
				}
			})
		}

		disabled := []string{"", "no", "false", "0", "disabled"}
		for _, val := range disabled {
			t.Run("disabled-"+val, func(t *testing.T) {
				os.Setenv(c.Name, val)
				if c.Func() == c.Expects {
					t.Fatalf("%s was enabled with falsy value: %s", name, val)
				}
			})
		}
	}
}

func TestEnvironmentMapEnabled(t *testing.T) {
	trueString := strconv.FormatBool(true)
	os.Setenv(env.ForceColor, trueString)
	os.Setenv(env.Debug, trueString)
	os.Setenv(env.Verbose, trueString)

	res := EnvironmentMap()
	if res == nil {
		t.Fatalf("Expected map, got nil")
	}

	expected := map[string]string{
		env.ForceColor: "always",
		env.Debug:      trueString,
		env.Verbose:    trueString,
	}

	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("Unexpected result from enabled environment. Wanted %v, got %v", res, expected)
	}
}

func TestEnvironmentMapDisabled(t *testing.T) {
	trueString := strconv.FormatBool(true)
	// clear COLOR
	os.Unsetenv(env.ForceColor)
	// set NO_COLOR
	os.Setenv(env.NoColor, trueString)
	os.Unsetenv(env.Debug)
	os.Unsetenv(env.Verbose)
	os.Setenv(env.Silent, trueString)

	res := EnvironmentMap()
	if res == nil {
		t.Fatalf("Expected map, got nil")
	}

	expected := map[string]string{
		env.NoColor: trueString,
		env.Silent:  trueString,
	}

	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("Unexpected result from disabled environment. Wanted %v, got %v", res, expected)
	}
}
