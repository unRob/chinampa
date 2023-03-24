// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package command

import (
	_c "git.rob.mx/nidito/chinampa/internal/constants"
	"git.rob.mx/nidito/chinampa/pkg/runtime"
)

var Root = &Command{
	Summary:     "Replace me with chinampa.Configure",
	Description: "",
	Path:        []string{runtime.Executable},
	Options: Options{
		_c.HelpCommandName: &Option{
			ShortName:   "h",
			Type:        "bool",
			Description: "Display help for any command",
		},
		"color": &Option{
			Type:        "bool",
			Description: "Always print colors to stderr",
			Default:     runtime.ColorEnabled(),
		},
		"no-color": &Option{
			Type:        "bool",
			Description: "Disable printing of colors to stderr",
			Default:     !runtime.ColorEnabled(),
		},
		"verbose": &Option{
			ShortName:   "v",
			Type:        "bool",
			Default:     runtime.VerboseEnabled(),
			Description: "Log verbose output to stderr",
		},
		"silent": &Option{
			Type:        "bool",
			Description: "Silence non-error logging",
		},
		"skip-validation": &Option{
			Type:        "bool",
			Description: "Do not validate any arguments or options",
		},
		"version": &Option{
			Type:        "bool",
			Default:     false,
			Description: "Display program version and exit",
		},
	},
}
