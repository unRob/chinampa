// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0

/*
package statuscode manages exit codes for programs.

See `man sysexits || grep "#define EX" /usr/include/sysexits.h`
and https://tldp.org/LDP/abs/html/exitcodes.html
*/
package statuscode

const (
	// Ok means everything is fine.
	Ok = 0
	// Usage means bad arguments were provided by the user.
	Usage = 64
	// ProgrammerError means the developer made a mistake.
	ProgrammerError = 70
	// ConfigError means configuration files or the environment is misconfigured.
	ConfigError = 78
	// NotFound means a sub-command not was found.
	NotFound = 127
)
