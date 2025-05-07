package commands

import (
	"Daemon/cmd/cli"
	"Daemon/internal/container"
	"context"
)

func init() {
	cli.RegisterCommand(&cli.Command{
		Name:        "remove",
		Description: "Remove a container by name",
		Args:        []cli.ArgSpec{{Name: "containerName", Required: true}},
		Execute:     runRemove,
	})
}

func runRemove(ctx context.Context, service *container.Service, args []string) error {
	name := args[0]

	if err := service.RemoveContainer(ctx, name); err != nil {
		return err
	}

	return nil
}
