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
	Execute     func(ctx context.Context, service *container.Service, args []string) error
}
