// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package command

import (
	"fmt"
	"strconv"
	"strings"

	"git.rob.mx/nidito/chinampa/pkg/errors"
	"git.rob.mx/nidito/chinampa/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Options is a map of name to Option.
type Options map[string]*Option

// AllKnown returns a map of option names to their resolved values.
func (opts *Options) AllKnown() map[string]any {
	col := map[string]any{}
	for name, opt := range *opts {
		col[name] = opt.ToValue()
	}
	return col
}

// AllKnownStr returns a map of option names to their stringified values.
func (opts *Options) AllKnownStr() map[string]string {
	col := map[string]string{}
	for name, opt := range *opts {
		col[name] = opt.ToString()
	}
	return col
}

// Parse populates values with those supplied in the provided pflag.Flagset.
func (opts *Options) Parse(supplied *pflag.FlagSet) {
	// log.Debugf("Parsing supplied flags, %v", supplied)
	for name, opt := range *opts {
		switch opt.Type {
		case ValueTypeBoolean:
			if val, err := supplied.GetBool(name); err == nil {
				opt.provided = val
				continue
			}
		case ValueTypeInt:
			if val, err := supplied.GetInt(name); err == nil {
				opt.provided = val
				continue
			}
		default:
			opt.Type = ValueTypeString
			if opt.Repeated {
				val, err := supplied.GetStringArray(name)
				if err == nil {
					opt.provided = val
					continue
				}
				logger.Errorf("Invalid option configuration: %s", err)
			} else {
				if val, err := supplied.GetString(name); err == nil {
					opt.provided = val
					continue
				}
			}
		}
	}
}

// AreValid tells if these options are all valid.
func (opts *Options) AreValid() error {
	for name, opt := range *opts {
		if err := opt.Validate(name); err != nil {
			return err
		}
	}

	return nil
}

// Option represents a command line flag.
type Option struct {
	// Type represents the type of value expected to be provided for this option.
	Type ValueType `json:"type" yaml:"type" validate:"omitempty,oneof=string bool int"`
	// Description is a required field that show up during completions and help.
	Description string `json:"description" yaml:"description" validate:"required"`
	// Default value for this option, if none provided.
	Default any `json:"default,omitempty" yaml:"default,omitempty"`
	// ShortName When set, enables representing this Option as a short flag (-x).
	ShortName string `json:"short-name,omitempty" yaml:"short-name,omitempty"` // nolint:tagliatelle
	// Values denote the source for completion/validation values of this option.
	Values *ValueSource `json:"values,omitempty" yaml:"values,omitempty" validate:"omitempty"`
	// Repeated options may be specified more than once.
	Repeated bool `json:"repeated" yaml:"repeated" validate:"omitempty"`
	// Command references the Command this Option is defined for.
	Command  *Command `json:"-" yaml:"-" validate:"-"`
	provided any
}

// IsKnown tells if the option was provided by the user.
func (opt *Option) IsKnown() bool {
	return opt.provided != nil
}

// Returns the resolved value for an option.
func (opt *Option) ToValue() any {
	if opt.IsKnown() {
		return opt.provided
	}
	return opt.Default
}

// Returns a string representation of this Option's resolved value.
func (opt *Option) ToString() string {
	value := opt.ToValue()
	stringValue := ""
	switch opt.Type {
	case ValueTypeBoolean:
		if value == nil {
			stringValue = ""
		} else {
			stringValue = strconv.FormatBool(value.(bool))
		}
	case ValueTypeInt:
		if value == nil {
			stringValue = ""
		} else {
			stringValue = fmt.Sprintf("%d", value)
		}
	default:
		if value != nil {
			stringValue = value.(string)
		}
	}

	return stringValue
}

func (opt *Option) internalValidate(name, current string) error {
	if current == "" {
		return nil
	}

	validValues, _, err := opt.Resolve(current)
	if err != nil {
		return err
	}

	if !contains(validValues, current) {
		return errors.BadArguments{Msg: fmt.Sprintf("%s is not a valid value for option <%s>. Valid options are: %s", current, name, strings.Join(validValues, ", "))}
	}

	return nil
}

// Validate validates the provided value if a value source.
func (opt *Option) Validate(name string) error {
	if !opt.Validates() {
		return nil
	}

	if opt.Repeated {
		values := opt.ToValue().([]string)
		for _, current := range values {
			if err := opt.internalValidate(name, current); err != nil {
				return err
			}
		}
	} else {
		if err := opt.internalValidate(name, opt.ToString()); err != nil {
			return err
		}
	}

	return nil
}

// Validates tells if the user-supplied value needs validation.
func (opt *Option) Validates() bool {
	return opt.Values != nil && opt.Values.Validates()
}

// providesAutocomplete tells if this option provides autocomplete values.
func (opt *Option) providesAutocomplete() bool {
	return opt.Values != nil
}

// Resolve returns autocomplete values for an option.
func (opt *Option) Resolve(currentValue string) (values []string, flag cobra.ShellCompDirective, err error) {
	if opt.Values != nil {
		if opt.Values.command == nil {
			opt.Values.command = opt.Command
		}
		return opt.Values.Resolve(currentValue)
	}

	return
}

// CompletionFunction is called by cobra when asked to complete an option.
func (opt *Option) CompletionFunction(cmd *cobra.Command, args []string, toComplete string) (values []string, flag cobra.ShellCompDirective) {
	if !opt.providesAutocomplete() {
		logger.Tracef("Option does not provide autocomplete %+v", opt)
		flag = cobra.ShellCompDirectiveNoFileComp
		return
	}

	if err := opt.Command.Arguments.Parse(args); err != nil {
		logger.Errorf("Could not parse command arguments %s", err)
		return []string{}, cobra.ShellCompDirectiveDefault
	}
	opt.Command.Options.Parse(cmd.Flags())

	var err error
	values, flag, err = opt.Resolve(toComplete)
	if err != nil {
		return values, cobra.ShellCompDirectiveError
	}

	if toComplete != "" && flag != cobra.ShellCompDirectiveFilterFileExt && flag != cobra.ShellCompDirectiveFilterDirs {
		filtered := []string{}
		for _, value := range values {
			if strings.HasPrefix(value, toComplete) {
				filtered = append(filtered, value)
			}
		}
		values = filtered
	}

	return cobra.AppendActiveHelp(values, opt.Description), flag
}
