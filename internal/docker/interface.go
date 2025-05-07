package docker

import "context"

type Client interface {
	CreateContainer(ctx context.Context, id, image string, env map[string]string, volumes []string, ports []int) (string, error)
	StartContainer(ctx context.Context, id string) error
	StopContainer(ctx context.Context, id string) error
	RemoveContainer(ctx context.Context, id string) error
}
