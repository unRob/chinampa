// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package command_test

import (
	"reflect"
	"strings"
	"testing"

	. "git.rob.mx/nidito/chinampa/pkg/command"
	"github.com/spf13/cobra"
)

func testCommand() *Command {
	return (&Command{
		Arguments: []*Argument{
			{
				Name:    "first",
				Default: "default",
			},
			{
				Name:     "variadic",
				Default:  []any{"defaultVariadic0", "defaultVariadic1"},
				Variadic: true,
			},
		},
		Options: Options{
			"option": {
				Default: "default",
				Type:    "string",
			},
			"bool": {
				Type:    "bool",
				Default: false,
			},
		},
	}).SetBindings()
}

func TestParse(t *testing.T) {
	cmd := testCommand()
	cmd.Arguments.Parse([]string{"asdf", "one", "two", "three"}) // nolint: errcheck
	known := cmd.Arguments.AllKnown()

	if !cmd.Arguments[0].IsKnown() {
		t.Fatalf("first argument isn't known")
	}
	val, exists := known["first"]
	if !exists {
		t.Fatalf("first argument isn't on AllKnown map: %v", known)
	}

	if val != "asdf" {
		t.Fatalf("first argument does not match. expected: %s, got %s", "asdf", val)
	}

	if !cmd.Arguments[1].IsKnown() {
		t.Fatalf("variadic argument isn't known")
	}
	val, exists = known["variadic"]
	if !exists {
		t.Fatalf("variadic argument isn't on AllKnown map: %v", known)
	}

	if !reflect.DeepEqual(val, []string{"one", "two", "three"}) {
		t.Fatalf("Known argument does not match. expected: %s, got %s", "one two three", val)
	}

	cmd = testCommand()
	cmd.Arguments.Parse([]string{"asdf"}) // nolint: errcheck
	known = cmd.Arguments.AllKnown()

	if !cmd.Arguments[0].IsKnown() {
		t.Fatalf("first argument is not known")
	}

	val, exists = known["first"]
	if !exists {
		t.Fatalf("first argument isn't on AllKnown map: %v", known)
	}

	if val != "asdf" {
		t.Fatalf("first argument does not match. expected: %s, got %s", "asdf", val)
	}

	val, exists = known["variadic"]
	if !exists {
		t.Fatalf("variadic argument isn't on AllKnown map: %v", known)
	}

	expected := []string{"defaultVariadic0", "defaultVariadic1"}
	if !reflect.DeepEqual(val, expected) {
		t.Fatalf("variadic argument does not match. expected: %s, got %s", expected, val)
	}
}

func TestBeforeParse(t *testing.T) {
	cmd := testCommand()
	known := cmd.Arguments.AllKnown()

	if cmd.Arguments[0].IsKnown() {
		t.Fatalf("first argument is known")
	}

	val, exists := known["first"]
	if !exists {
		t.Fatalf("first argument isn't on AllKnown map: %v", known)
	}

	if val != "default" {
		t.Fatalf("first argument does not match. expected: %s, got %s", "asdf", val)
	}

	val, exists = known["variadic"]
	if !exists {
		t.Fatalf("variadic argument isn't on AllKnown map: %v", known)
	}

	expected := []string{"defaultVariadic0", "defaultVariadic1"}
	if !reflect.DeepEqual(val, expected) {
		t.Fatalf("variadic argument does not match. expected: %s, got %s", expected, val)
	}
}

func TestArgumentsValidate(t *testing.T) {
	staticArgument := func(name string, def string, values []string, variadic bool) *Argument {
		return &Argument{
			Name:     name,
			Default:  def,
			Variadic: variadic,
			Required: def == "",
			Values: &ValueSource{
				Static: &values,
			},
		}
	}

	cases := []struct {
		Command     *Command
		Args        []string
		ErrorSuffix string
		Env         []string
	}{
		{
			Command: (&Command{
				// Name: []string{"test", "required", "failure"},
				Arguments: []*Argument{
					{
						Name:     "first",
						Required: true,
					},
				},
			}).SetBindings(),
			ErrorSuffix: "Missing argument for FIRST",
		},
		{
			Args:        []string{"bad"},
			ErrorSuffix: "bad is not a valid value for argument <first>. Valid options are: good, default",
			Command: (&Command{
				// Name: []string{"test", "script", "bad"},
				Arguments: []*Argument{
					{
						Name:    "first",
						Default: "default",
						Values: &ValueSource{
							Script: "echo good; echo default",
						},
					},
				},
			}).SetBindings(),
		},
		{
			Args:        []string{"bad"},
			ErrorSuffix: "bad is not a valid value for argument <first>. Valid options are: default, good",
			Command: (&Command{
				// Name: []string{"test", "static", "errors"},
				Arguments: []*Argument{staticArgument("first", "default", []string{"default", "good"}, false)},
			}).SetBindings(),
		},
		{
			Args:        []string{"default", "good", "bad"},
			ErrorSuffix: "bad is not a valid value for argument <first>. Valid options are: default, good",
			Command: (&Command{
				// Name:      []string{"test", "static", "errors"},
				Arguments: []*Argument{staticArgument("first", "default", []string{"default", "good"}, true)},
			}).SetBindings(),
		},
		{
			Args:        []string{"good"},
			ErrorSuffix: "could not validate argument for command test script bad-exit, ran",
			Command: (&Command{
				Path: []string{"test", "script", "bad-exit"},
				Arguments: []*Argument{
					{
						Name:    "first",
						Default: "default",
						Values: &ValueSource{
							Script: "echo good; echo default; exit 2",
						},
					},
				},
			}).SetBindings(),
		},
	}

	t.Run("good command is good", func(t *testing.T) {
		cmd := testCommand()
		cmd.Arguments[0] = staticArgument("first", "default", []string{"default", "good"}, false)
		cmd.Arguments[1] = staticArgument("second", "", []string{"one", "two", "three"}, true)
		cmd.SetBindings()

		cmd.Arguments.Parse([]string{"first", "one", "three", "two"}) // nolint: errcheck

		err := cmd.Arguments.AreValid()
		if err == nil {
			t.Fatalf("Unexpected failure validating: %s", err)
		}
	})

	for _, c := range cases {
		t.Run(c.Command.FullName(), func(t *testing.T) {
			c.Command.Arguments.Parse(c.Args) // nolint: errcheck

			err := c.Command.Arguments.AreValid()
			if err == nil {
				t.Fatalf("Expected failure but got none")
			}
			if !strings.HasPrefix(err.Error(), c.ErrorSuffix) {
				t.Fatalf("Could not find error <%s> got <%s>", c.ErrorSuffix, err)
			}
		})
	}
}

func TestArgumentToDesc(t *testing.T) {
	cases := []struct {
		Arg  *Argument
		Spec string
	}{
		{
			Arg: &Argument{
				Name: "regular",
			},
			Spec: "[REGULAR]",
		},
		{
			Arg: &Argument{
				Name:     "required",
				Required: true,
			},
			Spec: "REQUIRED",
		},
		{
			Arg: &Argument{
				Name:     "variadic-regular",
				Variadic: true,
			},
			Spec: "[VARIADIC_REGULAR...]",
		},
		{
			Arg: &Argument{
				Name:     "variadic-required",
				Variadic: true,
				Required: true,
			},
			Spec: "VARIADIC_REQUIRED...",
		},
	}

	for _, c := range cases {
		t.Run(c.Arg.Name, func(t *testing.T) {
			res := c.Arg.ToDesc()
			if res != c.Spec {
				t.Fatalf("Expected %s got %s", c.Spec, res)
			}
		})
	}
}

func TestArgumentCompletion(t *testing.T) {
	t.Parallel()

	testcmd := func() (*Command, *cobra.Command) {
		cc := &cobra.Command{
			Use:   "test-command",
			Short: "test",
		}

		cmd := testCommand()
		cmd.SetCobra(cc)
		cmd.SetBindings()
		return cmd, cc
	}

	t.Run("empty-args", func(t *testing.T) {
		t.Parallel()
		args := &Arguments{}

		cc := &cobra.Command{}
		values, directive := args.CompletionFunction(cc, []string{}, "")

		if len(values) != 0 {
			t.Fatal("Values offered for empty argument spec")
		}

		if directive != cobra.ShellCompDirectiveError {
			t.Fatalf("Unexpected directive: %d", directive)
		}
	})

	t.Run("empty arg spec", func(t *testing.T) {
		t.Parallel()
		cmd, cc := testcmd()
		values, directive := cmd.Arguments.CompletionFunction(cc, []string{}, "")

		if len(values) != 1 || values[0] != "_activeHelp_ " {
			t.Fatalf("Values offered for empty argument spec: (%d) %v", len(values), values)
		}

		if directive != cobra.ShellCompDirectiveError {
			t.Fatalf("Unexpected directive: %d", directive)
		}
	})

	t.Run("first-arg", func(t *testing.T) {
		t.Parallel()
		cmd, cc := testcmd()
		choices := []string{"au", "to", "com", "plete"}
		cmd.Arguments[0].Values = &ValueSource{
			Static: &choices,
		}

		values, directive := cmd.Arguments.CompletionFunction(cc, []string{}, "")

		expected := append(choices, "_activeHelp_") // nolint: gocritic
		if reflect.DeepEqual(expected, values) {
			t.Fatalf("Unexpected values offered: %v", values)
		}

		if directive != cobra.ShellCompDirectiveDefault {
			t.Fatalf("Unexpected directive: %d", directive)
		}
	})

	t.Run("first-arg-prefix", func(t *testing.T) {
		t.Parallel()
		choices := []string{"au", "to", "com", "plete"}
		cmd, cc := testcmd()
		cmd.Arguments[0].Values = &ValueSource{
			Static: &choices,
		}
		values, directive := cmd.Arguments.CompletionFunction(cc, []string{}, "a")

		expected := []string{"au", "_activeHelp_"}
		if reflect.DeepEqual(expected, values) {
			t.Fatalf("Unexpected values offered: %v", values)
		}

		if directive != cobra.ShellCompDirectiveDefault {
			t.Fatalf("Unexpected directive: %d", directive)
		}
	})

	t.Run("variadic-arg", func(t *testing.T) {
		choices := []string{"au", "to", "com", "plete"}
		cmd, cc := testcmd()
		cmd.Arguments[1].Values = &ValueSource{
			Static: &choices,
		}

		values, directive := cmd.Arguments.CompletionFunction(cc, []string{"au", ""}, "")

		expected := append(choices, "_activeHelp_") // nolint: gocritic
		if reflect.DeepEqual(expected, values) {
			t.Fatalf("Unexpected values offered: %v", values)
		}

		if directive != cobra.ShellCompDirectiveDefault {
			t.Fatalf("Unexpected directive: %d", directive)
		}
	})

	t.Run("variadic-arg-repeated", func(t *testing.T) {
		t.Parallel()
		cmd, cc := testcmd()
		choices := []string{"au", "to", "com", "plete"}
		cmd.Arguments[1].Values = &ValueSource{
			Static: &choices,
		}

		values, directive := cmd.Arguments.CompletionFunction(cc, []string{"au", "au", ""}, "t")

		expected := []string{"to", "_activeHelp_"}
		if reflect.DeepEqual(expected, values) {
			t.Fatalf("Unexpected values offered: %v", values)
		}

		if directive != cobra.ShellCompDirectiveDefault {
			t.Fatalf("Unexpected directive: %d", directive)
		}
	})
}
