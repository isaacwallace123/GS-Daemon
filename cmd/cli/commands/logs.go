package commands

import (
	"fmt"
	"io"
	"os"
	"os/signal"

	"golang.org/x/term"

	dockerContainer "github.com/docker/docker/api/types/container"

	"Daemon/cmd/cli"
	"Daemon/internal/docker"
	"Daemon/internal/shared/logger"
)

func init() {
	cli.RegisterCommand(&cli.Command{
		Name:        "logs",
		Description: "Attach to container's console (live logs + command input)",
		Args: []cli.ArgSpec{
			{Name: "containerName", Required: true},
		},
		Execute: runLogs,
	})
}

func runLogs(command *cli.CommandContext) error {
	name := command.Args[0]

	// Create Docker client
	client, err := docker.NewDockerClient()
	if err != nil {
		return logger.Error("Failed to initialize Docker client: %v", err)
	}

	// Resolve Docker container ID
	id, err := client.ResolveNameToID(command.Ctx, name)
	if err != nil {
		return logger.Error("Container not found: %v", err)
	}

	// Attach to the container
	attachResp, err := client.GetClient().ContainerAttach(command.Ctx, id, dockerContainer.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Logs:   true,
	})
	if err != nil {
		return logger.Error("Failed to attach: %v", err)
	}
	defer attachResp.Close()

	// Check terminal
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return logger.Error("stdin is not a terminal — cannot enter raw mode")
	}

	// Enter raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return logger.Error("Failed to set terminal raw mode: %v", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Println("🖥️  Attached to container. Type commands below. Press Ctrl+C to exit.\n")

	// Signal handler
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	// Pipe input
	go func() {
		_, _ = io.Copy(attachResp.Conn, os.Stdin)
	}()

	// Pipe output directly (TTY containers send raw stdout)
	go func() {
		_, _ = io.Copy(os.Stdout, attachResp.Reader)
	}()

	// Wait for Ctrl+C
	<-sigs
	fmt.Println("\n🛑 Exiting attach session.")

	return nil
}
