// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package render_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"git.rob.mx/nidito/chinampa/pkg/env"
	"git.rob.mx/nidito/chinampa/pkg/render"
)

func TestMarkdownUnstyled(t *testing.T) {
	content := []byte("# hello")
	os.Setenv(env.HelpStyle, "markdown")
	res, err := render.Markdown(content, false)

	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}

	expected := []byte("# hello") // nolint:ifshort
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("Unexpected response ---\n%s\n---\n wanted:\n---\n%s\n---", res, expected)
	}
}

func TestMarkdownNoColor(t *testing.T) {
	os.Unsetenv(env.HelpStyle)
	content := []byte("# hello ﹅world﹅")
	res, err := render.Markdown(content, false)

	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}

	// account for 80 character width word wrapping
	// our string is 15 characters, there's 2 for padding at the start
	spaces := "                                                             "

	expected := []byte("\n  # hello `world`" + spaces + "\n\n") // nolint:ifshort
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("Unexpected response ---\n%s\n---\n wanted:\n---\n%s\n---", res, expected)
	}
}

var autoStyleTestRender = "\n\x1b[38;5;228;48;5;63;1m\x1b[0m\x1b[38;5;228;48;5;63;1m\x1b[0m  \x1b[38;5;228;48;5;63;1m \x1b[0m\x1b[38;5;228;48;5;63;1mhello\x1b[0m\x1b[38;5;228;48;5;63;1m \x1b[0m\x1b[38;5;252m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[0m\n\x1b[0m\n"

const lightStyleTestRender = "\n\x1b[38;5;228;48;5;63;1m\x1b[0m\x1b[38;5;228;48;5;63;1m\x1b[0m  \x1b[38;5;228;48;5;63;1m \x1b[0m\x1b[38;5;228;48;5;63;1mhello\x1b[0m\x1b[38;5;228;48;5;63;1m \x1b[0m\x1b[38;5;234m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[0m\n\x1b[0m\n"

func TestMarkdownColor(t *testing.T) {
	os.Unsetenv(env.HelpStyle)
	content := []byte("# hello")

	styles := map[string][]byte{
		"":      []byte(autoStyleTestRender),
		"dark":  []byte(autoStyleTestRender),
		"auto":  []byte(autoStyleTestRender),
		"light": []byte(lightStyleTestRender),
	}
	for style, expected := range styles {
		t.Run(fmt.Sprintf("style %s", style), func(t *testing.T) {
			os.Setenv(env.HelpStyle, style)
			res, err := render.Markdown(content, true)

			if err != nil {
				t.Fatalf("Unexpected error %s", err)
			}

			if !reflect.DeepEqual(res, expected) {
				t.Fatalf("Unexpected response ---\n%v\n---\n wanted:\n---\n%v\n---", res, expected)
			}
		})
	}
}
