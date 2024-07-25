package shell

import (
	"log/slog"
	"os/exec"
)

type Command struct {
	Name string
	Args []string
}

func New(name string, args ...string) *Command {
	return &Command{
		Name: name,
		Args: args,
	}
}

func Start() chan Command {
	c := make(chan Command, 10)
	go func() {
		for cmd := range c {
			if cmd.Name == "stop" {
				break
			}

			slog.Info("new command", "name", cmd.Name, "args", cmd.Args)
			go func(cmd Command) {
				command := exec.Command(cmd.Name, cmd.Args...)

				err := command.Run()

				if err != nil {
					slog.Error("command failed", "err", err, "name", cmd.Name, "args", cmd.Args)
				}
			}(cmd)
		}
	}()

	return c
}
