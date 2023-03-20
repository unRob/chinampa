// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package registry

import (
	"fmt"
	"strings"

	"git.rob.mx/nidito/chinampa/pkg/command"
	"git.rob.mx/nidito/chinampa/pkg/errors"
	"git.rob.mx/nidito/chinampa/pkg/runtime"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func newCobraRoot(root *command.Command) *cobra.Command {
	return &cobra.Command{
		Use: root.Name() + " [--silent|-v|--verbose] [--[no-]color] [-h|--help] [--version]",
		Annotations: map[string]string{
			ContextKeyRuntimeIndex: root.Name(),
		},
		Short:             root.Summary,
		Long:              root.Description,
		DisableAutoGenTag: true,
		SilenceUsage:      true,
		SilenceErrors:     true,
		ValidArgs:         []string{""},
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.OnlyValidArgs(cmd, args); err != nil {
				suggestions := []string{}
				bold := color.New(color.Bold)
				for _, l := range cmd.SuggestionsFor(args[len(args)-1]) {
					suggestions = append(suggestions, bold.Sprint(l))
				}
				errMessage := fmt.Sprintf("Unknown subcommand %s", bold.Sprint(strings.Join(args, " ")))
				if len(suggestions) > 0 {
					errMessage += ". Perhaps you meant " + strings.Join(suggestions, ", ") + "?"
				}
				return errors.NotFound{Msg: errMessage, Group: []string{}}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				if ok, err := cmd.Flags().GetBool("version"); err == nil && ok {
					_, err := cmd.OutOrStdout().Write([]byte(cmd.Root().Annotations["version"]))
					return err
				}
				return errors.NotFound{Msg: "No subcommand provided", Group: []string{}}
			}
			return nil
		},
	}
}

func ToCobra(cmd *command.Command, globalOptions command.Options) *cobra.Command {
	localName := cmd.Name()
	useSpec := []string{localName, "[options]"}
	for _, arg := range cmd.Arguments {
		useSpec = append(useSpec, arg.ToDesc())
	}

	cc := &cobra.Command{
		Use:               strings.Join(useSpec, " "),
		Short:             cmd.Summary,
		DisableAutoGenTag: true,
		SilenceUsage:      true,
		SilenceErrors:     true,
		Hidden:            cmd.Hidden,
		Annotations: map[string]string{
			ContextKeyRuntimeIndex: cmd.FullName(),
		},
		Args: func(cc *cobra.Command, supplied []string) error {
			skipValidation, _ := cc.Flags().GetBool("skip-validation")
			if !skipValidation && runtime.ValidationEnabled() {
				cmd.Arguments.Parse(supplied)
				return cmd.Arguments.AreValid()
			}
			return nil
		},
		RunE: cmd.Run,
	}

	cc.SetFlagErrorFunc(func(c *cobra.Command, e error) error {
		return errors.BadArguments{Msg: e.Error()}
	})

	cc.ValidArgsFunction = cmd.Arguments.CompletionFunction

	cc.Flags().AddFlagSet(cmd.FlagSet())

	for name, opt := range cmd.Options {
		if err := cc.RegisterFlagCompletionFunc(name, opt.CompletionFunction); err != nil {
			log.Errorf("Failed setting up autocompletion for option <%s> of command <%s>", name, cmd.FullName())
		}
	}

	cc.SetHelpFunc(cmd.HelpRenderer(globalOptions))
	cmd.SetCobra(cc)
	return cc
}

func FromCobra(cc *cobra.Command) *command.Command {
	rtidx, hasAnnotation := cc.Annotations[ContextKeyRuntimeIndex]
	if hasAnnotation {
		return Get(rtidx)
	}
	return nil
}
