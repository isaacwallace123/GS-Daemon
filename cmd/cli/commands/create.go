package commands

import (
	"Daemon/cmd/cli"
	"Daemon/internal/shared/logger"
)

func init() {
	cli.RegisterCommand(&cli.Command{
		Name:        "create",
		Description: "Create a container from an egg",
		Args: []cli.ArgSpec{
			{Name: "eggName", Required: true},
			{Name: "containerName", Required: true},
		},
		Execute: runCreate,
	})
}

func runCreate(command *cli.CommandContext) error {
	eggName := command.Args[0]
	containerName := command.Args[1]

	_, err := command.Service.CreateContainer(command.Ctx, containerName, eggName)
	if err != nil {
		return err
	}

	logger.System("✅ Container '%s' created successfully.", containerName)

	return nil
}
