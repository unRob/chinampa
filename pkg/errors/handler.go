// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package errors

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	_c "git.rob.mx/nidito/chinampa/internal/constants"
	"git.rob.mx/nidito/chinampa/pkg/statuscode"
)

func showHelp(cmd *cobra.Command) {
	if cmd.Name() != _c.HelpCommandName {
		err := cmd.Help()
		if err != nil {
			os.Exit(statuscode.ProgrammerError)
		}
	}
}

// HandleCobraExit is called when a command errors out or was not found.
func HandleCobraExit(cmd *cobra.Command, err error) error {
	if err == nil {
		ok, err := cmd.PersistentFlags().GetBool(_c.HelpCommandName)
		if cmd.Name() == _c.HelpCommandName || err == nil && ok {
			os.Exit(statuscode.RenderHelp)
		}

		os.Exit(statuscode.Ok)
	}

	switch err.(type) {
	case BadArguments:
		showHelp(cmd)
		logrus.Error(err)
		os.Exit(statuscode.Usage)
	case NotFound:
		showHelp(cmd)
		logrus.Error(err)
		os.Exit(statuscode.NotFound)
	default:
		if strings.HasPrefix(err.Error(), "unknown command") {
			showHelp(cmd)
			os.Exit(statuscode.NotFound)
		} else if strings.HasPrefix(err.Error(), "unknown flag") || strings.HasPrefix(err.Error(), "unknown shorthand flag") {
			showHelp(cmd)
			logrus.Error(err)
			os.Exit(statuscode.Usage)
		}
	}

	logrus.Errorf("Unknown error: %s", err)
	os.Exit(2)
	return err
}
