package cli

import (
	"Daemon/internal/container"
	"Daemon/internal/shared/logger"
	"context"
)

var CommandMap = map[string]*Command{}

func RegisterCommand(cmd *Command) {
	CommandMap[cmd.Name] = cmd
}

func RunCommand(ctx context.Context, service *container.Service, input []string) {
	if len(input) == 0 {
		return
	}

	cmdName := input[0]
	args := input[1:]

	cmd, ok := CommandMap[cmdName]
	if !ok {
		logger.Info("❓ Unknown command. Type 'help' for a list.")
		return
	}

	// Validate args
	for i, spec := range cmd.Args {
		if spec.Required && i >= len(args) {
			logger.Error("❌ Missing required argument: <%s>", spec.Name)

			printCommandUsage(cmd)

			return
		}
	}

	// Run command
	if err := cmd.Execute(ctx, service, args); err != nil {
		logger.Error("❌ Command error: %v", err)
	}
}

func printCommandUsage(cmd *Command) {
	argHelp := ""

	for _, a := range cmd.Args {
		if a.Required {
			argHelp += " <" + a.Name + ">"
		} else {
			argHelp += " [" + a.Name + "]"
		}
	}

	logger.Info("Usage: %s%s", cmd.Name, argHelp)
}
