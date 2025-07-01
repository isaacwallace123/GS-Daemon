package commands

import (
	"Daemon/cmd/cli"
)

func init() {
	cli.RegisterCommand(&cli.Command{
		Name:        "logs",
		Description: "Print container logs",
		Args:        []cli.ArgSpec{{Name: "containerName", Required: true}},
		Execute:     runLogs,
	})
}

func runLogs(c *cli.CommandContext) error {
	name := c.Args[0]
	return c.Service.StreamLogs(c.Ctx, name)
}
