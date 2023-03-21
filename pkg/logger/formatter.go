// Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>
// SPDX-License-Identifier: Apache-2.0
package logger

import (
	"strings"
	"time"

	"git.rob.mx/nidito/chinampa/pkg/runtime"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var bold *color.Color
var boldRedBG *color.Color
var boldRed *color.Color
var boldYellowBG *color.Color
var boldYellow *color.Color
var dimmed *color.Color

func init() {
	bold = color.New(color.Bold)
	bold.EnableColor()
	boldRedBG = color.New(color.Bold, color.BgRed)
	boldRedBG.EnableColor()
	boldRed = color.New(color.Bold, color.FgHiRed)
	boldRed.EnableColor()
	boldYellowBG = color.New(color.Bold, color.BgYellow, color.FgBlack)
	boldYellowBG.EnableColor()
	boldYellow = color.New(color.Bold, color.FgHiYellow)
	boldYellow.EnableColor()
	dimmed = color.New(color.Faint)
	dimmed.EnableColor()
}

type Formatter struct {
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	prefix := ""
	colorEnabled := runtime.ColorEnabled()
	message := entry.Message
	switch {
	case runtime.VerboseEnabled():
		date := strings.Replace(entry.Time.Local().Format(time.DateTime), " ", "T", 1)
		component := ""
		if c, ok := entry.Data[componentKey]; ok {
			component = " " + c.(string)
		}
		level := entry.Level.String()
		if colorEnabled {
			switch {
			case entry.Level <= logrus.ErrorLevel:
				level = boldRed.Sprint(level)
			case entry.Level == logrus.WarnLevel:
				level = boldYellow.Sprint(level)
			case entry.Level >= logrus.DebugLevel:
				level = dimmed.Sprint(level)
				message = dimmed.Sprint(message)
			default:
				level = dimmed.Sprint(level)
			}
		}

		prefix = dimmed.Sprint(date) + " " + level + dimmed.Sprint(component) + "\t"
	case entry.Level == logrus.ErrorLevel:
		if colorEnabled {
			prefix = boldRedBG.Sprint(" ERROR ") + " "
		} else {
			prefix = "ERROR: "
		}
	case entry.Level == logrus.WarnLevel:
		if colorEnabled {
			prefix = boldYellowBG.Sprint(" WARNING ") + " "
		} else {
			prefix = "WARNING: "
		}
	case entry.Level >= logrus.DebugLevel:
		if colorEnabled {
			prefix = dimmed.Sprintf("%s: ", strings.ToUpper(entry.Level.String()))
			message = dimmed.Sprint(message)
		} else {
			prefix = strings.ToUpper(entry.Level.String()) + ": "
		}
	}

	return []byte(prefix + message + "\n"), nil
}
