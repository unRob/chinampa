// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package command

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type varSearchMap struct {
	Status int
	Name   string
	Usage  string
}

func (cmd *Command) Validate() (report map[string]int) {
	report = map[string]int{}

	validate := validator.New()
	if err := validate.Struct(cmd); err != nil {
		verrs := err.(validator.ValidationErrors)
		for _, issue := range verrs {
			// todo: output better errors, see validator.FieldError
			report[fmt.Sprint(issue)] = 1
		}
	}

	vars := map[string]map[string]*varSearchMap{
		"argument": {},
		"option":   {},
	}

	for _, arg := range cmd.Arguments {
		vars["argument"][strings.ToUpper(strings.ReplaceAll(arg.Name, "-", "_"))] = &varSearchMap{2, arg.Name, ""}
	}

	for name := range cmd.Options {
		vars["option"][strings.ToUpper(strings.ReplaceAll(name, "-", "_"))] = &varSearchMap{2, name, ""}
	}

	return report
}
