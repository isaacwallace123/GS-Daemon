package commands

import (
	"Daemon/cmd/cli"
	"Daemon/internal/shared/logger"
)

func init() {
	cli.RegisterCommand(&cli.Command{
		Name:        "stop",
		Description: "Stop a container by name",
		Args:        []cli.ArgSpec{{Name: "containerName", Required: true}},
		Execute:     runStop,
	})
}

func runStop(c *cli.CommandContext) error {
	name := c.Args[0]

	if err := c.Service.StopContainer(c.Ctx, name); err != nil {
		return err
	}

	logger.System("✅ Container '%s' stopped successfully.", name)

	return nil
}
