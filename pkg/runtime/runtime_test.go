// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package runtime_test

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"git.rob.mx/nidito/chinampa/pkg/env"
	. "git.rob.mx/nidito/chinampa/pkg/runtime"
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
		ResetParsedFlagsCache()

		for k := range env {
			os.Unsetenv(k)
		}
		for _, entry := range prevEnv {
			parts := strings.SplitN(entry, "=", 2)
			os.Setenv(parts[0], parts[1])
		}
	})
}

func TestCombinations(t *testing.T) {
	args := append([]string{}, os.Args...)
	t.Cleanup(func() { os.Args = args })
	cases := []struct {
		Env     map[string]string
		Args    []string
		Func    func() bool
		Expects bool
	}{
		{
			Env:     map[string]string{},
			Args:    []string{},
			Func:    VerboseEnabled,
			Expects: false,
		},
		{
			Env:     map[string]string{env.Verbose: "1"},
			Args:    []string{"--silent"},
			Func:    VerboseEnabled,
			Expects: false,
		},
		{
			Env:     map[string]string{env.Verbose: "1"},
			Args:    []string{},
			Func:    VerboseEnabled,
			Expects: true,
		},
		{
			Env:     map[string]string{env.Silent: "1"},
			Args:    []string{},
			Func:    VerboseEnabled,
			Expects: false,
		},
		{
			Env:     map[string]string{},
			Args:    []string{},
			Func:    SilenceEnabled,
			Expects: false,
		},
		{
			Env:     map[string]string{env.Silent: "1"},
			Args:    []string{},
			Func:    SilenceEnabled,
			Expects: true,
		},
		{
			Env:     map[string]string{env.Silent: "1"},
			Args:    []string{"--verbose"},
			Func:    SilenceEnabled,
			Expects: false,
		},
		{
			Env:     map[string]string{env.Verbose: "1"},
			Args:    []string{"--silent"},
			Func:    SilenceEnabled,
			Expects: true,
		},
		{
			Env:     map[string]string{env.ForceColor: "1"},
			Args:    []string{"--no-color"},
			Func:    ColorEnabled,
			Expects: false,
		},
		{
			Env:     map[string]string{},
			Args:    []string{},
			Func:    ColorEnabled,
			Expects: true,
		},
		{
			Env:     map[string]string{env.ForceColor: "1"},
			Args:    []string{},
			Func:    ColorEnabled,
			Expects: true,
		},
		{
			Env:     map[string]string{env.ForceColor: "1"},
			Args:    []string{"--no-color"},
			Func:    ColorEnabled,
			Expects: false,
		},
		{
			Env:     map[string]string{env.NoColor: "1"},
			Args:    []string{},
			Func:    ColorEnabled,
			Expects: false,
		},
		{
			Env:     map[string]string{env.NoColor: "1"},
			Args:    []string{"--color"},
			Func:    ColorEnabled,
			Expects: true,
		},
	}

	for _, c := range cases {
		fname := runtime.FuncForPC(reflect.ValueOf(c.Func).Pointer()).Name()
		name := fmt.Sprintf("%v/%v/%s", fname, c.Env, c.Args)
		t.Run(name, func(t *testing.T) {
			withEnv(t, c.Env)
			os.Args = c.Args
			if res := c.Func(); res != c.Expects {
				t.Fatalf("%s got %v wanted: %v", name, res, c.Expects)
			}
		})
	}
}

func TestEnabled(t *testing.T) {
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
			Name:    env.Silent,
			Func:    SilenceEnabled,
			Expects: true,
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
				withEnv(t, map[string]string{c.Name: val})
				if c.Func() != c.Expects {
					t.Fatalf("%s wasn't enabled with a valid value: %s", name, val)
				}
			})
		}

		disabled := []string{"", "no", "false", "0", "disabled"}
		for _, val := range disabled {
			t.Run("disabled-"+val, func(t *testing.T) {
				withEnv(t, map[string]string{c.Name: val})
				if c.Func() == c.Expects {
					t.Fatalf("%s was enabled with falsy value: %s", name, val)
				}
			})
		}
	}
}

func TestSilent(t *testing.T) {
	args := append([]string{}, os.Args...)
	t.Cleanup(func() { os.Args = args })
	t.Run("SILENT = silence", func(t *testing.T) {
		withEnv(t, map[string]string{
			env.Silent:  "1",
			env.Verbose: "",
		})
		os.Args = []string{}
		if !SilenceEnabled() {
			t.Fail()
		}
	})

	t.Run("SILENT+VERBOSE=silence", func(t *testing.T) {
		withEnv(t, map[string]string{
			env.Silent:  "1",
			env.Verbose: "1",
		})
		os.Args = []string{}
		if !SilenceEnabled() {
			t.Fail()
		}
	})

	t.Run("VERBOSE+--silent=silent", func(t *testing.T) {
		withEnv(t, map[string]string{
			env.Silent:  "0",
			env.Verbose: "1",
		})
		os.Args = []string{"some", "random", "--silent", "args"}
		if !SilenceEnabled() {
			t.Fail()
		}
	})

	t.Run("--silent=silent", func(t *testing.T) {
		withEnv(t, map[string]string{
			env.Silent:  "",
			env.Verbose: "",
		})
		os.Args = []string{"some", "random", "--silent", "args"}
		if !SilenceEnabled() {
			t.Fail()
		}
	})

	t.Run("nothing = nothing", func(t *testing.T) {
		withEnv(t, map[string]string{
			env.Silent:  "",
			env.Verbose: "",
		})
		os.Args = []string{"some", "random", "args"}
		if SilenceEnabled() {
			t.Fail()
		}
		if VerboseEnabled() {
			t.Fail()
		}
	})
}

func TestEnvironmentMapEnabled(t *testing.T) {
	trueString := strconv.FormatBool(true)
	withEnv(t, map[string]string{
		env.ForceColor: trueString,
		env.Debug:      trueString,
		env.Verbose:    trueString,
	})

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
