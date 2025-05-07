package commands

import (
	"Daemon/cmd/cli"
	"Daemon/internal/container"
	"context"
)

func init() {
	cli.RegisterCommand(&cli.Command{
		Name:        "stop",
		Description: "Stop a container by name",
		Args:        []cli.ArgSpec{{Name: "containerName", Required: true}},
		Execute:     runStop,
	})
}

func runStop(ctx context.Context, service *container.Service, args []string) error {
	name := args[0]

	if err := service.StopContainer(ctx, name); err != nil {
		return err
	}

	return nil
}
