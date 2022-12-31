// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package chinampa

import (
	"git.rob.mx/nidito/chinampa/internal/registry"
	"git.rob.mx/nidito/chinampa/pkg/command"
	"git.rob.mx/nidito/chinampa/pkg/runtime"
)

func Register(cmd *command.Command) {
	registry.Register(cmd.SetBindings())
}

type Config struct {
	Name        string
	Version     string
	Summary     string
	Description string
}

func Execute(config Config) error {
	runtime.Executable = config.Name
	command.Root.Summary = config.Summary
	command.Root.Description = config.Description
	command.Root.Path = []string{runtime.Executable}
	return registry.Execute(config.Version)
}
