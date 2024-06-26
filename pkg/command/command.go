// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package command

import (
	"fmt"
	"strconv"
	"strings"

	"git.rob.mx/nidito/chinampa/pkg/logger"
	"git.rob.mx/nidito/chinampa/pkg/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var log = logger.Sub("chinampa:command")

type HelpFunc func(printLinks bool) string
type Action func(cmd *Command) error

type Command struct {
	Path []string `json:"path" yaml:"path"`
	// Summary is a short description of a command, on supported shells this is part of the autocomplete prompt
	Summary string `json:"summary" yaml:"summary" validate:"required"`
	// Description is a long form explanation of how a command works its magic. Markdown is supported
	Description string `json:"description" yaml:"description" validate:"required"`
	// A list of arguments for a command
	Arguments Arguments `json:"arguments" yaml:"arguments" validate:"dive"`
	// A map of option names to option definitions
	Options  Options  `json:"options" yaml:"options" validate:"dive"`
	HelpFunc HelpFunc `json:"-" yaml:"-"`
	// The action to take upon running
	Action       Action `json:"-" yaml:"-"`
	runtimeFlags *pflag.FlagSet
	Cobra        *cobra.Command `json:"-" yaml:"-"`
	// Meta stores application specific stuff
	Meta   any  `json:"meta" yaml:"meta"`
	Hidden bool `json:"-" yaml:"-"`
}

func (cmd *Command) IsRoot() bool {
	return cmd.FullName() == runtime.Executable
}

func (cmd *Command) SetBindings() *Command {
	ptr := cmd
	for _, opt := range cmd.Options {
		opt.Command = ptr
		if opt.Validates() {
			opt.Values.command = ptr
		}
	}

	for _, arg := range cmd.Arguments {
		arg.Command = ptr
		if arg.Validates() {
			arg.Values.command = ptr
		}
	}
	return ptr
}

func (cmd *Command) Name() string {
	return cmd.Path[len(cmd.Path)-1]
}

func (cmd *Command) FullName() string {
	return strings.Join(cmd.Path, " ")
}

func (cmd *Command) FlagSet() *pflag.FlagSet {
	if cmd.runtimeFlags == nil {
		fs := pflag.NewFlagSet(strings.Join(cmd.Path, " "), pflag.ContinueOnError)
		fs.SortFlags = false
		fs.Usage = func() {}
		for name, opt := range cmd.Options {
			switch opt.Type {
			case ValueTypeBoolean:
				def := false
				if opt.Default != nil {
					def = opt.Default.(bool)
				}
				fs.BoolP(name, opt.ShortName, def, opt.Description)
			case ValueTypeInt:
				def := -1
				if opt.Default != nil {
					switch val := opt.Default.(type) {
					case int:
						def = opt.Default.(int)
					case string:
						casted, err := strconv.Atoi(val)
						if err != nil {
							log.Warnf("Could not parse default with value <%s> as integer for option <%s>", val, name)
						}
						def = casted
					}
				}
				fs.IntP(name, opt.ShortName, def, opt.Description)
			case ValueTypeDefault, ValueTypeString:
				opt.Type = ValueTypeString
				if opt.Repeated {
					def := []string{}
					if opt.Default != nil {
						switch defV := opt.Default.(type) {
						case []any:
							for _, v := range defV {
								def = append(def, fmt.Sprintf("%s", v))
							}
						case []string:
							def = defV
						case string:
							def = []string{defV}
						default:
							logger.Errorf("Invalid default for repeated option %s configuration: %+v", name, defV)
						}
					}

					fs.StringArrayP(name, opt.ShortName, def, opt.Description)
				} else {
					def := ""
					if opt.Default != nil {
						def = fmt.Sprintf("%s", opt.Default)
					}
					fs.StringP(name, opt.ShortName, def, opt.Description)
				}
			default:
				// ignore flag
				log.Warnf("Ignoring unknown option type <%s> for option <%s>", opt.Type, name)
				continue
			}
		}

		cmd.runtimeFlags = fs
	}
	return cmd.runtimeFlags
}

func (cmd *Command) ParseInput(cc *cobra.Command, args []string) error {
	if err := cmd.Arguments.Parse(args); err != nil {
		return err
	}
	skipValidation, _ := cc.Flags().GetBool("skip-validation")
	cmd.Options.Parse(cc.Flags())
	if !skipValidation {
		log.Debug("Validating arguments")
		if err := cmd.Arguments.AreValid(); err != nil {
			return err
		}

		log.Debug("Validating flags")
		if err := cmd.Options.AreValid(); err != nil {
			log.Debugf("Invalid flags for %s: %s", cmd.FullName(), err)
			return err
		}
	}

	return nil
}

func (cmd *Command) Run(cc *cobra.Command, args []string) error {
	log.Debugf("running command %s", cmd.FullName())

	if err := cmd.ParseInput(cc, args); err != nil {
		log.Debugf("Parsing input to command %s failed: %s", cmd.FullName(), err)
		return err
	}

	return cmd.Action(cmd)
}

func (cmd *Command) SetCobra(cc *cobra.Command) {
	cmd.Cobra = cc
}
