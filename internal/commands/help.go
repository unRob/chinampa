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

func cobraCommandFullName(c *cobra.Command) []string {
	name := []string{c.Name()}
	if c.HasParent() {
		name = append(cobraCommandFullName(c.Parent()), name...)
	}
	return name
}

var Help = &cobra.Command{
	Use:   _c.HelpCommandName + " [command]",
	Short: "Display usage information for any command",
	Long: `Help provides the valid arguments and options for any command known to ﹅@chinampa@﹅.

## Colorized output

By default, and unless ﹅` + env.NoColor + `﹅ is set, ﹅@chinampa@ help﹅ will query the environment variable ﹅COLORFGBG﹅ to decide which style to use when rendering help, unless if ﹅` + env.HelpStyle + `﹅ is set to any of the following values: **light**, **dark**, **markdown**, and **auto**. 24-bit color is available when ﹅COLORTERM﹅ is set to ﹅truecolor﹅.`,
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
		c.Long = strings.ReplaceAll(c.Long, "@chinampa@", runtime.Executable)
		if len(args) > 0 && c != nil && c.Name() != args[len(args)-1] {
			c, topicArgs, err := c.Root().Find(args)
			if err == nil && c != nil && len(topicArgs) == 0 {
				// exact command help
				cobra.CheckErr(c.Help())
				os.Exit(statuscode.RenderHelp)
				return
			}

			if err != nil {
				logrus.Error(err)
				os.Exit(statuscode.ProgrammerError)
			}

			fullName := strings.Join(cobraCommandFullName(c), " ")
			cobra.CheckErr(c.Help())
			if len(topicArgs) > 0 {
				logrus.Errorf("Unknown help topic \"%s\" for %s", topicArgs[0], fullName)
			} else {
				logrus.Errorf("Unknown help topic \"%s\"", args[0])
			}
			os.Exit(statuscode.NotFound)
		} else {
			c.InitDefaultHelpFlag() // make possible 'help' flag to be shown
			cobra.CheckErr(c.Help())
		}

		os.Exit(statuscode.RenderHelp)
	},
}
