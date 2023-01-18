# Chinampa

Like [`milpa`](https://milpa.dev) but for go programs only.

```go
package main

import (
	"os"
	"fmt"

	"git.rob.mx/nidito/chinampa"
	"git.rob.mx/nidito/chinampa/pkg/command"
	"git.rob.mx/nidito/chinampa/pkg/runtime"
	"github.com/sirupsen/logrus"
)


func main() {
	chinampa.Register(&command.Command{
		Path:        []string{"something"},
		Summary:     "does something",
		Description: "a longer description of how it does stuff",
		Arguments: command.Arguments{
			{
				Name:        "argument zero",
				Description: "a help text for using argument zero",
				Required:    true,
			},
		},
		Action: func(cmd *command.Command) error {
			someArg := cmd.Arguments[0].ToValue().(string)

			return fmt.Errorf("Don't know how to do stuff with %s", someArg)
		},
	})

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableLevelTruncation: true,
		DisableTimestamp:       true,
		ForceColors:            runtime.ColorEnabled(),
	})

	if runtime.DebugEnabled() {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Debugging enabled")
	}

	cfg := chinampa.Config{
		Name:        "myProgram",
		Version:     "0.0.0",
		Summary:     "a short summary of my program",
		Description: "a longer text describing what its for",
	}

	if err := chinampa.Execute(cfg); err != nil {
		logrus.Errorf("total failure: %s", err)
		os.Exit(2)
	}
}
```
