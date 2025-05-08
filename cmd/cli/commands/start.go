package commands

import (
	"Daemon/cmd/cli"
	"Daemon/internal/shared/logger"
)

func init() {
	cli.RegisterCommand(&cli.Command{
		Name:        "start",
		Description: "Start a container by name",
		Args:        []cli.ArgSpec{{Name: "containerName", Required: true}},
		Execute:     runStart,
	})
}

func runStart(c *cli.CommandContext) error {
	name := c.Args[0]

	if err := c.Service.StartContainer(c.Ctx, name); err != nil {
		return err
	}

	logger.System("✅ Container '%s' started successfully.", name)

	return nil
}
