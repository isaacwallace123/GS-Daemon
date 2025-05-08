package commands

import (
	"Daemon/cmd/cli"
	"Daemon/internal/shared/logger"
)

func init() {
	cli.RegisterCommand(&cli.Command{
		Name:        "get",
		Description: "Get container details by name",
		Args:        []cli.ArgSpec{{Name: "containerName", Required: true}},
		Execute:     runGet,
	})
}

func runGet(command *cli.CommandContext) error {
	name := command.Args[0]

	container, err := command.Service.GetContainer(command.Ctx, name)
	if err != nil {
		return err
	}

	logger.Info("Retrieved container: %+v", container)

	return nil
}
