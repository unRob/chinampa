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
