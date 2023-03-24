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
	colorEnabled := runtime.ColorEnabled()
	c := &combinedCommand{
		Spec:          cmd,
		Command:       cmd.Cobra,
		GlobalOptions: globalOptions,
		HTMLOutput:    runtime.HelpStyle() == "markdown",
	}
	err := render.HelpTemplate(runtime.Executable).Execute(&buf, c)
	if err != nil {
		return nil, err
	}

	content, err := render.Markdown(buf.Bytes(), colorEnabled)
	if err != nil {
		return nil, err
	}
	return content, nil
}
