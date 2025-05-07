package cli

import (
	"Daemon/internal/container"
	"Daemon/internal/shared/logger"
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

func RunShell(ctx context.Context, service *container.Service) {
	logger.Info("Welcome to the container CLI. Type 'help' to see commands.")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		args := strings.Fields(line)

		if len(args) == 0 {
			continue
		}

		if args[0] == "exit" || args[0] == "quit" {
			logger.Info("👋 Exiting...")
			return
		}

		if args[0] == "help" {
			printHelp()
			continue
		}

		RunCommand(ctx, service, args)
	}
}

func printHelp() {
	var builder strings.Builder

	builder.WriteString("📖 Available commands:\n")

	for _, cmd := range CommandMap {
		argHelp := ""

		for _, arg := range cmd.Args {
			if arg.Required {
				argHelp += " <" + arg.Name + ">"
			} else {
				argHelp += " [" + arg.Name + "]"
			}
		}

		builder.WriteString("  " + cmd.Name + argHelp + " - " + cmd.Description + "\n")
	}

	builder.WriteString("  help              - Show this help message\n")
	builder.WriteString("  exit / quit       - Exit CLI\n")

	logger.Info("\n%s", builder.String())
}
