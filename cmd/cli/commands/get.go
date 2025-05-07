package commands

import (
	"Daemon/cmd/cli"
	"Daemon/internal/container"
	"Daemon/internal/shared/logger"
	"context"
)

func init() {
	cli.RegisterCommand(&cli.Command{
		Name:        "get",
		Description: "Get container details by name",
		Args:        []cli.ArgSpec{{Name: "containerName", Required: true}},
		Execute:     runGet,
	})
}

func runGet(ctx context.Context, service *container.Service, args []string) error {
	name := args[0]
	c, err := service.GetContainer(ctx, name)

	if err != nil {
		return err
	}

	logger.Info("Retrieved container: %+v", c)

	return nil
}
