// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0

/*
Package env holds environment variable names that are meant to be overridden by implementations.

# example

	package main

	import "git.rob.mx/nidito/chinampa/env"

	func init() {
		env.HelpUnstyled = "MY_APP_PLAIN_HELP"
		env.HelpStyle = "MY_APP_HELP_STYLE"
		env.Verbose = "MY_APP_VERBOSE"
		env.Silent = "MY_APP_SILENT"
		env.ValidationDisabled = "MY_APP_SKIP_VALIDATION"
	}
*/
package env

// HelpUnstyled means help will not be colored nor formatted for a TTY.
var HelpUnstyled = "HELP_STYLE_PLAIN"

// HelpStyle identifies the theme to use for help formatting.
var HelpStyle = "HELP_STYLE"

// Verbose enables verbose printing of log entries.
var Verbose = "VERBOSE"

// Silent disables all printing of log entries, except for errors.
var Silent = "SILENT"

// NoColor disables printing of color escape codes in help and log entries.
var NoColor = "NO_COLOR"

// ForceColor enables printing of color escape codes in help and log entries.
var ForceColor = "COLOR"

// ValidationDisabled disables validation on arguments and options.
var ValidationDisabled = "SKIP_VALIDATION"

// Debug enables printing of debugging information.
var Debug = "DEBUG"
