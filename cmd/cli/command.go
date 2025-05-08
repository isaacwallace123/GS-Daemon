package cli

import (
	"context"

	"Daemon/internal/container"
)

type ArgSpec struct {
	Name     string
	Required bool
}

type Command struct {
	Name        string
	Description string
	Args        []ArgSpec
	Execute     func(c *CommandContext) error
}

type CommandContext struct {
	Ctx     context.Context
	Service *container.Service
	Args    []string
}
