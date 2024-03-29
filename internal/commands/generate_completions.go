// SPDX-License-Identifier: Apache-2.0
// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
package commands

import (
	"os"

	"github.com/spf13/cobra"
)

var GenerateCompletions = &cobra.Command{
	Use:               "__generate_completions [bash|zsh|fish]",
	Short:             "Outputs a shell-specific script for autocompletions that can be piped into a file",
	Hidden:            true,
	DisableAutoGenTag: true,
	SilenceUsage:      true,
	Args:              cobra.ExactArgs(1),
	ValidArgs:         []string{"bash", "fish", "zsh"},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		switch args[0] {
		case "bash":
			err = cmd.Root().GenBashCompletionV2(os.Stdout, true)
		case "zsh":
			err = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			err = cmd.Root().GenFishCompletion(os.Stdout, true)
		}
		return
	},
}
