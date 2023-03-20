package logger

import (
	"fmt"
	"strings"
	"time"

	"git.rob.mx/nidito/chinampa/pkg/runtime"
	"github.com/sirupsen/logrus"
)

type Formatter struct {
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	prefix := ""
	colorEnabled := runtime.ColorEnabled()
	message := entry.Message
	if runtime.VerboseEnabled() {
		date := strings.Replace(entry.Time.Local().Format(time.DateTime), " ", "T", 1)
		component := ""
		if c, ok := entry.Data[componentKey]; ok {
			component = " " + c.(string)
		}
		level := entry.Level.String()
		if colorEnabled {
			if entry.Level <= logrus.ErrorLevel {
				level = "\033[31m\033[1m" + level + "\033[0m"
			} else if entry.Level == logrus.WarnLevel {
				level = "\033[33m\033[1m" + level + "\033[0m"
			} else if entry.Level >= logrus.DebugLevel && colorEnabled {
				message = "\033[2m" + message + "\033[0m"
			}
		}

		prefix = fmt.Sprintf("\033[2m%s %s%s\033[0m\t", date, level, component)
	} else if entry.Level == logrus.ErrorLevel {
		if colorEnabled {
			prefix = "\033[41m\033[1m ERROR \033[0m "
		} else {
			prefix = "ERROR: "
		}
	} else if entry.Level == logrus.WarnLevel {
		if colorEnabled {
			prefix = "\033[43m\033[31m warning \033[0m "
			message = "\033[33m" + message + "\033[0m"
		} else {
			prefix = "WARNING: "
		}
	} else if entry.Level >= logrus.DebugLevel {
		if colorEnabled {
			prefix = "\033[2m" + entry.Level.String() + ":\033[0m "
			message = "\033[2m" + message + "\033[0m"
		} else {
			prefix = strings.ToUpper(entry.Level.String()) + ": "
		}
	}
	return []byte(prefix + message + "\n"), nil
}
