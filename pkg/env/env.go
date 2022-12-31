// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package env

// Environment Variables.
var HelpUnstyled = "HELP_STYLE_PLAIN"
var HelpStyle = "HELP_STYLE"
var Verbose = "VERBOSE"
var Silent = "SILENT"
var NoColor = "NO_COLOR"
var ForceColor = "COLOR"
var ValidationDisabled = "SKIP_VALIDATION"
var Debug = "DEBUG"

// FlagNames are flags also available as environment variables.
var FlagNames = map[string]string{
	"no-color":        NoColor,
	"color":           ForceColor,
	"silent":          Silent,
	"verbose":         Verbose,
	"skip-validation": ValidationDisabled,
}
