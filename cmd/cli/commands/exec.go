package commands

import (
	"Daemon/cmd/cli"
	"Daemon/internal/shared/logger"
)

func init() {
	cli.RegisterCommand(&cli.Command{
		Name:        "exec",
		Description: "Run a command inside a container",
		Args: []cli.ArgSpec{
			{Name: "containerName", Required: true},
			{Name: "cmd", Required: true},
		},
		Execute: runExec,
	})
}

func runExec(c *cli.CommandContext) error {
	if len(c.Args) < 2 {
		return logger.Error("usage: exec <containerName> <command>")
	}
	name := c.Args[0]
	cmd := c.Args[1:]
	return c.Service.ExecCommand(c.Ctx, name, cmd)
}
