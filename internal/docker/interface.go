package docker

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"io"
)

type Client interface {
	CreateContainer(
		ctx context.Context,
		id, image string,
		env map[string]string,
		volumes []string,
		ports []int,
		name string,
	) (string, error)
	StartContainer(ctx context.Context, name string) error
	StopContainer(ctx context.Context, name string) error
	RemoveContainer(ctx context.Context, name string) error
	GetContainerByID(ctx context.Context, id string) (*container.InspectResponse, error)

	ResolveNameToID(ctx context.Context, name string) (string, error)

	// ExecInteractive runs a command inside the given container ID and
	// attaches the current terminal for interactive input/output.
	ExecInteractive(ctx context.Context, id string, cmd []string) error

	// GetContainerLogs returns a reader for the container's logs.
	GetContainerLogs(ctx context.Context, id string) (io.ReadCloser, error)
}
