// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package registry

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"git.rob.mx/nidito/chinampa/internal/commands"
	"git.rob.mx/nidito/chinampa/pkg/command"
	"git.rob.mx/nidito/chinampa/pkg/errors"
	"git.rob.mx/nidito/chinampa/pkg/statuscode"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ContextKeyRuntimeIndex is the string key used to store context in a cobra Command.
const ContextKeyRuntimeIndex = "x-chinampa-runtime-index"

var registry = &CommandRegistry{
	kv: map[string]*command.Command{},
}

type ByPath []*command.Command

func (cmds ByPath) Len() int           { return len(cmds) }
func (cmds ByPath) Swap(i, j int)      { cmds[i], cmds[j] = cmds[j], cmds[i] }
func (cmds ByPath) Less(i, j int) bool { return cmds[i].FullName() < cmds[j].FullName() }

type CommandRegistry struct {
	kv     map[string]*command.Command
	byPath []*command.Command
}

func Register(cmd *command.Command) {
	logrus.Debugf("Registering %s", cmd.FullName())
	registry.kv[cmd.FullName()] = cmd
}

func Get(id string) *command.Command {
	return registry.kv[id]
}

func CommandList() []*command.Command {
	if len(registry.byPath) == 0 {
		list := []*command.Command{}
		for _, v := range registry.kv {
			list = append(list, v)
		}
		sort.Sort(ByPath(list))
		registry.byPath = list
	}

	return registry.byPath
}

func Execute(version string) error {
	cmdRoot := command.Root
	ccRoot := newCobraRoot(command.Root)
	ccRoot.Annotations["version"] = version
	ccRoot.CompletionOptions.HiddenDefaultCmd = true
	ccRoot.PersistentFlags().AddFlagSet(cmdRoot.FlagSet())
	ccRoot.SetHelpCommand(commands.Help)
	ccRoot.AddCommand(commands.Version)
	ccRoot.AddCommand(commands.GenerateCompletions)

	for name, opt := range cmdRoot.Options {
		if err := ccRoot.RegisterFlagCompletionFunc(name, opt.CompletionFunction); err != nil {
			logrus.Errorf("Failed setting up autocompletion for option <%s> of command <%s>", name, cmdRoot.FullName())
		}
	}
	// ccRoot.SetHelpFunc(func(cc *cobra.Command, args []string) {
	// 	cmdRoot.HelpRenderer(cmdRoot.Options)(cc, args)
	// 	os.Exit(statuscode.RenderHelp)
	// })
	ccRoot.SetHelpFunc(cmdRoot.HelpRenderer(cmdRoot.Options))

	for _, cmd := range CommandList() {
		cmd := cmd
		container := ccRoot
		for idx, cp := range cmd.Path {
			if idx == len(cmd.Path)-1 {
				leaf := ToCobra(cmd, cmdRoot.Options)
				logrus.Debugf("adding command %s to %s", leaf.Name(), container.Name())
				container.AddCommand(leaf)
				break
			}

			query := []string{cp}
			found := false
			if cp == "help" && container == ccRoot {
				container = commands.Help
			} else {
				for _, sub := range container.Commands() {
					if sub.Name() == cp {
						container = sub
						found = true
					}
				}
			}

			if !found {
				groupName := strings.Join(query, " ")
				groupPath := append(cmdRoot.Path, append(cmd.Path[0:idx], query...)...) // nolint:gocritic
				cc := &cobra.Command{
					Use:                        cp,
					Short:                      fmt.Sprintf("%s subcommands", groupName),
					DisableAutoGenTag:          true,
					SuggestionsMinimumDistance: 2,
					SilenceUsage:               true,
					SilenceErrors:              true,
					Annotations: map[string]string{
						ContextKeyRuntimeIndex: strings.Join(groupPath, " "),
					},
					Args: func(cmd *cobra.Command, args []string) error {
						if err := cobra.OnlyValidArgs(cmd, args); err == nil {
							return nil
						}

						suggestions := []string{}
						bold := color.New(color.Bold)
						for _, l := range cmd.SuggestionsFor(args[len(args)-1]) {
							suggestions = append(suggestions, bold.Sprint(l))
						}
						last := len(args) - 1
						parent := cmd.CommandPath()
						errMessage := fmt.Sprintf("Unknown subcommand %s of known command %s", bold.Sprint(args[last]), bold.Sprint(parent))
						if len(suggestions) > 0 {
							errMessage += ". Perhaps you meant " + strings.Join(suggestions, ", ") + "?"
						}
						return errors.NotFound{Msg: errMessage, Group: []string{}}
					},
					ValidArgs: []string{""},
					RunE: func(cc *cobra.Command, args []string) error {
						if len(args) == 0 {
							return errors.NotFound{Msg: "No subcommand provided", Group: []string{}}
						}
						os.Exit(statuscode.NotFound)
						return nil
					},
				}

				groupParent := &command.Command{
					Path:        groupPath,
					Summary:     fmt.Sprintf("%s subcommands", groupName),
					Description: fmt.Sprintf("Runs subcommands within %s", groupName),
					Arguments:   command.Arguments{},
					Options:     command.Options{},
				}
				Register(groupParent)
				cc.SetHelpFunc(groupParent.HelpRenderer(command.Options{}))
				container.AddCommand(cc)
				container = cc
			}
		}

		cmd.Path = append(cmdRoot.Path, cmd.Path...)
	}
	cmdRoot.SetCobra(ccRoot)

	current, remaining, err := ccRoot.Find(os.Args[1:])
	if err != nil {
		current = ccRoot
	}
	logrus.Debugf("Chinampa found command %s, remaining %s", current.Name(), remaining)

	// if current.HasSubCommands() && current.
	if sub, _, err := current.Find(remaining); err == nil && sub != current {
		logrus.Debugf("Chinampa found sub-command %s, of %s", sub.Name(), current.Name())
		current = sub
	}
	logrus.Debugf("Chinampa is going to call command %s", current.Name())
	err = current.Execute()
	if err != nil {
		logrus.Debugf("Chinampa found error calling command %s", current.Name())
		errors.HandleCobraExit(current, err)
	}

	return err
}
