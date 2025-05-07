package commands

import (
	"context"

	"Daemon/cmd/cli"
	"Daemon/internal/container"
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

func runCreate(ctx context.Context, service *container.Service, args []string) error {
	eggName := args[0]
	containerName := args[1]

	_, err := service.CreateContainer(ctx, containerName, eggName)

	if err != nil {
		return err
	}

	logger.System("✅ Container '%s' created successfully.", containerName)

	return nil
}
