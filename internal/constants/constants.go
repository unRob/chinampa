// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package constants

import (
	// Embed requires an import so the compiler knows what's up. Golint requires a comment. Gotta please em both.
	_ "embed"
)

// HelpCommandName sets the name for the command that offers help.
const HelpCommandName = "help"

// HelpTemplate is the markdown template to use when rendering help.
//
//go:embed help.md
var HelpTemplate string
