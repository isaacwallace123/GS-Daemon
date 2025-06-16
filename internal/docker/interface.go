package docker

import (
	"context"
	"github.com/docker/docker/api/types/container"
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
}
