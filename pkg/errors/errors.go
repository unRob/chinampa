// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package errors

type NotFound struct {
	Msg   string
	Group []string
}

type BadArguments struct {
	Msg string
}

func (err NotFound) Error() string {
	return err.Msg
}

func (err BadArguments) Error() string {
	return err.Msg
}
