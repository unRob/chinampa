// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package errors

// NotFound happens when a sub-command was not found.
type NotFound struct {
	Msg   string
	Group []string
}

// BadArguments happens when the user provided incorrect arguments or options.
type BadArguments struct {
	Msg string
}

func (err NotFound) Error() string {
	return err.Msg
}

func (err BadArguments) Error() string {
	return err.Msg
}
