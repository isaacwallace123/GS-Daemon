package commands

import (
	"time"

	"Daemon/cmd/cli"
	"Daemon/internal/shared/logger"
)

func init() {
	cli.RegisterCommand(&cli.Command{
		Name:        "demo",
		Description: "Run a 30-second paper demo container",
		Args:        []cli.ArgSpec{},
		Execute:     runDemo,
	})
}

func runDemo(command *cli.CommandContext) error {
	logger.Info("Running 30-second container demo...")

	newContainer, err := command.Service.CreateContainer(command.Ctx, "paper-demo", "paper")
	if err != nil {
		return err
	}

	if err := command.Service.StartContainer(command.Ctx, newContainer.Name); err != nil {
		return err
	}
	time.Sleep(5 * time.Second)

	if err := command.Service.StopContainer(command.Ctx, newContainer.Name); err != nil {
		return err
	}

	if err := command.Service.RemoveContainer(command.Ctx, newContainer.Name); err != nil {
		return err
	}

	logger.System("✅ Demo complete.")
	return nil
}
