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
	"git.rob.mx/nidito/chinampa/pkg/logger"
	"github.com/charmbracelet/glamour"
	"github.com/sirupsen/logrus"
	"golang.org/x/term"
)

func addBackticks(str []byte) []byte {
	return bytes.ReplaceAll(str, []byte("﹅"), []byte("`"))
}

// Markdown renders markdown-formatted content to the tty.
func Markdown(content []byte, withColor bool) ([]byte, error) {
	content = addBackticks(content)
	var styleFunc glamour.TermRendererOption

	style := os.Getenv(env.HelpStyle)
	if style == "markdown" {
		// markdown will render frontmatter along content and will not format for
		// tty readability
		return content, nil
	}
	if withColor {
		switch style {
		case "dark":
			// For color TTYs with light text on dark background
			styleFunc = glamour.WithStandardStyle("dark")
		case "light":
			// For color TTYs with dark text on light background
			styleFunc = glamour.WithStandardStyle("light")
		default:
			// Glamour selects a style for the user.
			styleFunc = glamour.WithStandardStyle("auto")
			if style != "" {
				logger.Warnf("Unknown %s=%s, assuming \"auto\"", env.HelpStyle, style)
			}
		}
	} else {
		// basically the same as the "markdown" style, except formatted for
		// tty redability, prettifying and indenting, while not adding color.
		styleFunc = glamour.WithStandardStyle("notty")
	}

	width, _, err := term.GetSize(0)
	if err != nil {
		logrus.Debugf("Could not get terminal width")
		width = 80
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
