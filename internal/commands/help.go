// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package commands

import (
	"fmt"
	"os"
	"strings"

	_c "git.rob.mx/nidito/chinampa/internal/constants"
	"git.rob.mx/nidito/chinampa/pkg/env"
	"git.rob.mx/nidito/chinampa/pkg/runtime"
	"git.rob.mx/nidito/chinampa/pkg/statuscode"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Help = &cobra.Command{
	Use:   _c.HelpCommandName + " [command]",
	Short: "Display usage information for any command",
	Long:  `Help provides the valid arguments and options for any command known to ` + runtime.Executable + `. By default, ﹅` + runtime.Executable + ` help﹅ will query the environment variable ﹅COLORFGBG﹅ to decide which style to use when rendering help, except if ﹅` + env.HelpUnstyled + `﹅ is set. Valid styles are: **light**, **dark**, and **auto**.`,
	ValidArgsFunction: func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var completions []string
		cmd, _, e := c.Root().Find(args)
		if e != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		if cmd == nil {
			// Root help command.
			cmd = c.Root()
		}
		for _, subCmd := range cmd.Commands() {
			if subCmd.IsAvailableCommand() || subCmd.Name() == _c.HelpCommandName {
				if strings.HasPrefix(subCmd.Name(), toComplete) {
					completions = append(completions, fmt.Sprintf("%s\t%s", subCmd.Name(), subCmd.Short))
				}
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(c *cobra.Command, args []string) {
		cmd, _, e := c.Root().Find(args)
		if cmd == nil || e != nil || (len(args) > 0 && cmd != nil && cmd.Name() != args[len(args)-1]) {
			if cmd == nil {
				err := c.Root().Help()
				if err != nil {
					logrus.Error(err)
					os.Exit(statuscode.ProgrammerError)
				}
				logrus.Errorf("Unknown help topic %s", args)
				os.Exit(statuscode.NotFound)
			} else {
				err := cmd.Help()

				if err != nil {
					logrus.Error(err)
					os.Exit(statuscode.ProgrammerError)
				}

				if len(args) > 1 {
					logrus.Errorf("Unknown help topic %s for %s", args[1], args[0])
				} else {
					logrus.Errorf("Unknown help topic %s for %s", runtime.Executable, args[0])
				}
				os.Exit(statuscode.NotFound)
			}
		} else {
			cmd.InitDefaultHelpFlag() // make possible 'help' flag to be shown
			cobra.CheckErr(cmd.Help())
		}

		os.Exit(statuscode.RenderHelp)
	},
}
