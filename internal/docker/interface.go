package docker

import "context"

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

	ResolveNameToID(ctx context.Context, name string) (string, error)
}
