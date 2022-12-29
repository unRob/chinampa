// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
		"verbose": &Option{
			ShortName:   "v",
			Type:        "bool",
			Default:     runtime.VerboseEnabled(),
			Description: "Log verbose output to stderr",
		},
		"version": &Option{
			Type:        "bool",
			Default:     false,
			Description: "Display program version and exit",
		},
		"no-color": &Option{
			Type:        "bool",
			Description: "Disable printing of colors to stderr",
			Default:     !runtime.ColorEnabled(),
		},
		"color": &Option{
			Type:        "bool",
			Description: "Always print colors to stderr",
			Default:     runtime.ColorEnabled(),
		},
		"silent": &Option{
			Type:        "bool",
			Description: "Silence non-error logging",
		},
		"skip-validation": &Option{
			Type:        "bool",
			Description: "Do not validate any arguments or options",
		},
	},
}
