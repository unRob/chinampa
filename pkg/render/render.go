// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package render

import (
	"bytes"
	"os"
	"strings"
	"text/template"

	_c "git.rob.mx/nidito/chinampa/internal/constants"
	"git.rob.mx/nidito/chinampa/pkg/env"
	"git.rob.mx/nidito/chinampa/pkg/runtime"
	"github.com/charmbracelet/glamour"
	"github.com/sirupsen/logrus"
	"golang.org/x/term"
)

func addBackticks(str []byte) []byte {
	return bytes.ReplaceAll(str, []byte("﹅"), []byte("`"))
}

func Markdown(content []byte, withColor bool) ([]byte, error) {
	content = addBackticks(content)

	if runtime.UnstyledHelpEnabled() {
		return content, nil
	}

	width, _, err := term.GetSize(0)
	if err != nil {
		logrus.Debugf("Could not get terminal width")
		width = 80
	}

	var styleFunc glamour.TermRendererOption

	if withColor {
		style := os.Getenv(env.HelpStyle)
		switch style {
		case "dark":
			styleFunc = glamour.WithStandardStyle("dark")
		case "light":
			styleFunc = glamour.WithStandardStyle("light")
		default:
			styleFunc = glamour.WithStandardStyle("auto")
		}
	} else {
		styleFunc = glamour.WithStandardStyle("notty")
	}

	renderer, err := glamour.NewTermRenderer(
		styleFunc,
		glamour.WithEmoji(),
		glamour.WithWordWrap(width),
	)

	if err != nil {
		return content, err
	}

	return renderer.RenderBytes(content)
}

// TemplateFuncs is a FuncMap with aliases to the strings package.
var TemplateFuncs = template.FuncMap{
	"contains":   strings.Contains,
	"hasSuffix":  strings.HasSuffix,
	"hasPrefix":  strings.HasPrefix,
	"replace":    strings.ReplaceAll,
	"toUpper":    strings.ToUpper,
	"toLower":    strings.ToLower,
	"trim":       strings.TrimSpace,
	"trimSuffix": strings.TrimSuffix,
	"trimPrefix": strings.TrimPrefix,
	"list":       func(args ...any) []any { return args },
}

// TemplateCommandHelp holds a template for rendering command help.
var TemplateCommandHelp *template.Template

func HelpTemplate(executableName string) *template.Template {
	if TemplateCommandHelp == nil {
		TemplateCommandHelp = template.Must(template.New("help").Funcs(TemplateFuncs).Parse(strings.ReplaceAll(_c.HelpTemplate, "@chinampa@", executableName)))
	}

	return TemplateCommandHelp
}
