package commands

import (
	"Daemon/cmd/cli"
	"Daemon/internal/shared/logger"
)

func init() {
	cli.RegisterCommand(&cli.Command{
		Name:        "remove",
		Description: "Remove a container by name",
		Args:        []cli.ArgSpec{{Name: "containerName", Required: true}},
		Execute:     runRemove,
	})
}

func runRemove(c *cli.CommandContext) error {
	name := c.Args[0]

	if err := c.Service.RemoveContainer(c.Ctx, name); err != nil {
		return err
	}

	logger.System("✅ Container '%s' removed successfully.", name)

	return nil
}
