// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package statuscode

// Exit statuses
// see man sysexits || grep "#define EX" /usr/include/sysexits.h
// and https://tldp.org/LDP/abs/html/exitcodes.html
const (
	// 0 means everything is fine.
	Ok = 0
	// 42 provides answers to life, the universe and everything; also, renders help.
	RenderHelp = 42
	// 64 bad arguments
	// EX_USAGE The command was used incorrectly, e.g., with the wrong number of arguments, a bad flag, a bad syntax in a parameter, or whatever.
	Usage = 64
	// EX_SOFTWARE An internal software error has been detected. This should be limited to non-operating system related errors as possible.
	ProgrammerError = 70
	// EX_CONFIG Something was found in an unconfigured or misconfigured state.
	ConfigError = 78
	// 127 command not found.
	NotFound = 127
)
