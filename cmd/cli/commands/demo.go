package commands

import (
	"context"
	"time"

	"Daemon/cmd/cli"
	"Daemon/internal/container"
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

func runDemo(ctx context.Context, service *container.Service, _ []string) error {
	logger.Info("Running 30-second container demo...")

	newContainer, err := service.CreateContainer(ctx, "paper-demo", "paper")

	if err != nil {
		return err
	}

	err = service.StartContainer(ctx, newContainer.Name)

	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	err = service.StopContainer(ctx, newContainer.Name)

	if err != nil {
		return err
	}

	err = service.RemoveContainer(ctx, newContainer.Name)

	if err != nil {
		return err
	}

	logger.System("✅ Demo complete.")

	return nil
}
