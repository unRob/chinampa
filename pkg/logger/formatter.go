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
	if runtime.VerboseEnabled() {
		date := strings.Replace(entry.Time.Local().Format(time.DateTime), " ", "T", 1)
		component := ""
		if c, ok := entry.Data[componentKey]; ok {
			component = " " + c.(string)
		}
		prefix = fmt.Sprintf("\033[2m%s %s%s\033[0m\t", date, entry.Level.String(), component)
	} else if entry.Level == logrus.ErrorLevel {
		if colorEnabled {
			prefix = "\033[41m\033[1m ERROR \033[0m "
		} else {
			prefix = "ERROR: "
		}
	} else if entry.Level == logrus.WarnLevel {
		if colorEnabled {
			prefix = "\033[43m\033[31m warning \033[0m "
		} else {
			prefix = "WARNING: "
		}
	}
	return []byte(prefix + entry.Message + "\n"), nil
}
