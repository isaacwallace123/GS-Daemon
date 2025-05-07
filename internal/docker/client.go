package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"Daemon/internal/shared/logger"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	cli *client.Client
}

func NewDockerClient() (*DockerClient, error) {
	os.Setenv("DOCKER_API_VERSION", "1.44")

	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return nil, logger.Error("Failed to initialize Docker client: %v", err)
	}

	return &DockerClient{cli: cli}, nil
}

func (d *DockerClient) CreateContainer(
	ctx context.Context,
	id, imagePath string,
	env map[string]string,
	volumes []string,
	ports []int,
	name string,
) (string, error) {
	logger.Docker("Pulling image: %s", imagePath)

	reader, err := d.cli.ImagePull(ctx, imagePath, image.PullOptions{})

	if err != nil {
		return "", logger.Error("Failed to pull image: %v", err)
	}

	io.Copy(io.Discard, reader)

	var envList []string

	for key, val := range env {
		envList = append(envList, fmt.Sprintf("%s=%s", key, val))
	}

	resp, err := d.cli.ContainerCreate(ctx, &container.Config{
		Image: imagePath,
		Env:   envList,
		Labels: map[string]string{
			"com.daemon.uuid": id,
			"com.daemon.name": name,
			"created_by":      "wings-cli",
		},
	}, nil, nil, nil, name)
	if err != nil {
		return "", logger.Error("Failed to create container: %v", err)
	}

	logger.Docker("Container created: DockerID=%s", resp.ID)

	return resp.ID, nil
}

func (d *DockerClient) StartContainer(ctx context.Context, id string) error {
	logger.Docker("Starting container: %s", id)

	if err := d.cli.ContainerStart(ctx, id, container.StartOptions{}); err != nil {
		return logger.Error("Failed to start container %s: %v", id, err)
	}

	logger.Docker("Container started: %s", id)

	return nil
}

func (d *DockerClient) StopContainer(ctx context.Context, id string) error {
	logger.Docker("Stopping container: %s", id)

	opts := container.StopOptions{}

	if err := d.cli.ContainerStop(ctx, id, opts); err != nil {
		return logger.Error("Failed to stop container %s: %v", id, err)
	}

	logger.Docker("Container stopped: %s", id)

	return nil
}

func (d *DockerClient) RemoveContainer(ctx context.Context, id string) error {
	logger.Docker("Removing container: %s", id)

	if err := d.cli.ContainerRemove(ctx, id, container.RemoveOptions{Force: true}); err != nil {
		return logger.Error("Failed to remove container %s: %v", id, err)
	}

	logger.Docker("Container removed: %s", id)

	return nil
}

func (d *DockerClient) ResolveNameToID(ctx context.Context, name string) (string, error) {
	containers, err := d.cli.ContainerList(ctx, container.ListOptions{All: true})

	if err != nil {
		return "", err
	}

	for _, container := range containers {
		if container.Labels["com.daemon.name"] == name {
			return container.ID, nil
		}
	}

	return "", fmt.Errorf("container with name '%s' not found", name)
}
