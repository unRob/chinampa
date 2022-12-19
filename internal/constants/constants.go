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
package constants

import (
	"strings"
	"text/template"

	// Embed requires an import so the compiler knows what's up. Golint requires a comment. Gotta please em both.
	_ "embed"
)

const HelpCommandName = "help"

// Environment Variables.
const EnvVarHelpUnstyled = "HELP_STYLE_PLAIN"
const EnvVarHelpStyle = "HELP_STYLE"
const EnvVarMilpaVerbose = "VERBOSE"
const EnvVarMilpaSilent = "SILENT"
const EnvVarMilpaUnstyled = "NO_COLOR"
const EnvVarMilpaForceColor = "COLOR"
const EnvVarValidationDisabled = "SKIP_VALIDATION"
const EnvVarDebug = "DEBUG"

// EnvFlagNames are flags also available as environment variables.
var EnvFlagNames = map[string]string{
	"no-color":        EnvVarMilpaUnstyled,
	"color":           EnvVarMilpaForceColor,
	"silent":          EnvVarMilpaSilent,
	"verbose":         EnvVarMilpaVerbose,
	"skip-validation": EnvVarValidationDisabled,
}

// Exit statuses
// see man sysexits || grep "#define EX" /usr/include/sysexits.h
// and https://tldp.org/LDP/abs/html/exitcodes.html
const (
	// 0 means everything is fine.
	ExitStatusOk = 0
	// 42 provides answers to life, the universe and everything; also, renders help.
	ExitStatusRenderHelp = 42
	// 64 bad arguments
	// EX_USAGE The command was used incorrectly, e.g., with the wrong number of arguments, a bad flag, a bad syntax in a parameter, or whatever.
	ExitStatusUsage = 64
	// EX_SOFTWARE An internal software error has been detected. This should be limited to non-operating system related errors as possible.
	ExitStatusProgrammerError = 70
	// EX_CONFIG Something was found in an unconfigured or misconfigured state.
	ExitStatusConfigError = 78
	// 127 command not found.
	ExitStatusNotFound = 127
)

// ContextKeyRuntimeIndex is the string key used to store context in a cobra Command.
const ContextKeyRuntimeIndex = "x-chinampa-runtime-index"

//go:embed help.md
var helpTemplateText string

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
}

// TemplateCommandHelp holds a template for rendering command help.
var TemplateCommandHelp = template.Must(template.New("help").Funcs(TemplateFuncs).Parse(helpTemplateText))
