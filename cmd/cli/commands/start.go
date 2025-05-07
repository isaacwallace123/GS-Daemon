package commands

import (
	"Daemon/cmd/cli"
	"Daemon/internal/container"
	"context"
)

func init() {
	cli.RegisterCommand(&cli.Command{
		Name:        "start",
		Description: "Start a container by name",
		Args:        []cli.ArgSpec{{Name: "containerName", Required: true}},
		Execute:     runStart,
	})
}

func runStart(ctx context.Context, service *container.Service, args []string) error {
	name := args[0]

	if err := service.StartContainer(ctx, name); err != nil {
		return err
	}

	return nil
}
