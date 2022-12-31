// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package command

import (
	"bytes"

	"git.rob.mx/nidito/chinampa/pkg/render"
	"git.rob.mx/nidito/chinampa/pkg/runtime"
	"github.com/spf13/cobra"
)

type combinedCommand struct {
	Spec          *Command
	Command       *cobra.Command
	GlobalOptions Options
	HTMLOutput    bool
}

func (cmd *Command) HasAdditionalHelp() bool {
	return cmd.HelpFunc != nil
}

func (cmd *Command) AdditionalHelp(printLinks bool) *string {
	if cmd.HelpFunc != nil {
		str := cmd.HelpFunc(printLinks)
		return &str
	}
	return nil
}

func (cmd *Command) HelpRenderer(globalOptions Options) func(cc *cobra.Command, args []string) {
	return func(cc *cobra.Command, args []string) {
		// some commands don't have a binding until help is rendered
		// like virtual ones (sub command groups)
		cmd.SetCobra(cc)
		content, err := cmd.ShowHelp(globalOptions, args)
		if err != nil {
			panic(err)
		}
		_, err = cc.OutOrStderr().Write(content)
		if err != nil {
			panic(err)
		}
	}
}

func (cmd *Command) ShowHelp(globalOptions Options, args []string) ([]byte, error) {
	var buf bytes.Buffer
	c := &combinedCommand{
		Spec:          cmd,
		Command:       cmd.Cobra,
		GlobalOptions: globalOptions,
		HTMLOutput:    runtime.UnstyledHelpEnabled(),
	}
	err := render.HelpTemplate(runtime.Executable).Execute(&buf, c)
	if err != nil {
		return nil, err
	}

	colorEnabled := runtime.ColorEnabled()
	flags := cmd.Cobra.Flags()
	ncf := cmd.Cobra.Flag("no-color") // nolint:ifshort
	cf := cmd.Cobra.Flag("color")     // nolint:ifshort

	if noColorFlag, err := flags.GetBool("no-color"); err == nil && ncf.Changed {
		colorEnabled = !noColorFlag
	} else if colorFlag, err := flags.GetBool("color"); err == nil && cf.Changed {
		colorEnabled = colorFlag
	}

	content, err := render.Markdown(buf.Bytes(), colorEnabled)
	if err != nil {
		return nil, err
	}
	return content, nil
}
