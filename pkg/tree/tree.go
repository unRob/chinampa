// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package tree

import (
	"bytes"
	"sort"

	"git.rob.mx/nidito/chinampa/internal/registry"
	"git.rob.mx/nidito/chinampa/pkg/command"
	"github.com/spf13/cobra"
)

type CommandTree struct {
	Command  *command.Command `json:"command"`
	Children []*CommandTree   `json:"children"`
}

func (t *CommandTree) Traverse(fn func(cmd *command.Command) error) error {
	for _, child := range t.Children {
		if err := fn(child.Command); err != nil {
			return err
		}

		if err := child.Traverse(fn); err != nil {
			return err
		}
	}
	return nil
}

var tree *CommandTree

func Build(cc *cobra.Command, depth int) {
	root := registry.FromCobra(cc)
	if root == nil && cc.Root() == cc {
		root = command.Root
	}
	tree = &CommandTree{
		Command:  root,
		Children: []*CommandTree{},
	}

	var populateTree func(cmd *cobra.Command, ct *CommandTree, maxDepth int, depth int)
	populateTree = func(cmd *cobra.Command, ct *CommandTree, maxDepth int, depth int) {
		newDepth := depth + 1
		for _, subcc := range cmd.Commands() {
			if subcc.Hidden {
				continue
			}

			if cmd := registry.FromCobra(subcc); cmd != nil {
				leaf := &CommandTree{Children: []*CommandTree{}}
				leaf.Command = cmd
				ct.Children = append(ct.Children, leaf)

				if newDepth < maxDepth {
					populateTree(subcc, leaf, maxDepth, newDepth)
				}
			}
		}
	}
	populateTree(cc, tree, depth, 0)
}

func Serialize(serializationFn func(any) ([]byte, error)) (string, error) {
	content, err := serializationFn(tree)
	if err != nil {
		return "", err
	}
	return string(bytes.ReplaceAll(content, []byte("﹅"), []byte("`"))), nil
}

func ChildrenNames() []string {
	if tree == nil {
		return []string{}
	}

	ret := make([]string, len(tree.Children))
	for idx, cmd := range tree.Children {
		ret[idx] = cmd.Command.Name()
	}
	sort.Strings(ret)
	return ret
}

func CommandList() []*command.Command {
	return registry.CommandList()
}
