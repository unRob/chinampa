// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package commands

import (
	"os"

	"git.rob.mx/nidito/chinampa/pkg/statuscode"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Version = &cobra.Command{
	Use:               "version",
	Short:             "Display program version",
	Hidden:            false,
	DisableAutoGenTag: true,
	SilenceUsage:      true,
	RunE: func(cmd *cobra.Command, args []string) error {
		output := cmd.ErrOrStderr()
		version := cmd.Root().Annotations["version"]
		if cmd.CalledAs() == "" {
			// user asked for --version directly
			output = cmd.OutOrStderr()
			version += "\n"
		}

		_, err := output.Write([]byte(version))
		if err != nil {
			logrus.Errorf("version error: %s", err)
			return err
		}

		os.Exit(statuscode.Ok)
		return nil
	},
}
