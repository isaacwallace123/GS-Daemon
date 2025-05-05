package service

import (
	"Daemon/internal/core"
	"Daemon/pkg/logger"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type DockerDeployer struct {
	client *client.Client
	ctx    context.Context
}

func NewDockerDeployer() (*DockerDeployer, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)

	if err != nil {
		return nil, err
	}

	return &DockerDeployer{client: cli, ctx: context.Background()}, nil
}

func (d *DockerDeployer) Run(egg core.Egg) error {
	logger.Deployer("Deploying container for: %s", egg.GetName())

	_, _, err := d.client.ImageInspectWithRaw(d.ctx, egg.GetImage())
	if err != nil {
		logger.Deployer("Image %s not found locally. Pulling...", egg.GetImage())

		out, err := d.client.ImagePull(d.ctx, egg.GetImage(), image.PullOptions{})

		if err != nil {
			return logger.Error("Failed to pull image: %w", err)
		}

		defer out.Close()

		io.Copy(io.Discard, out)
	}

	logger.Debug("Final environment variables: %v", formatEnvVars(egg.GetEnvironment()))

	portBindings, exposedPorts := convertPorts(egg.GetPorts())
	mounts := convertVolumes(egg.GetVolumes())

	resp, err := d.client.ContainerCreate(d.ctx, &container.Config{
		Image:        egg.GetImage(),
		Cmd:          []string{"/bin/sh", "-c", egg.GetStartupCommand()},
		ExposedPorts: exposedPorts,
		Env:          formatEnvVars(egg.GetEnvironment()),
	}, &container.HostConfig{
		PortBindings: portBindings,
		Mounts:       mounts,
	}, &network.NetworkingConfig{}, nil, egg.GetName())

	if err != nil {
		return logger.Error("Container creation failed: %w", err)
	}

	if err := d.client.ContainerStart(d.ctx, resp.ID, container.StartOptions{}); err != nil {
		return logger.Error("Failed to start container: %w", err)
	}

	logger.Deployer("✅ Container %s started for %s", resp.ID[:12], egg.GetName())

	return nil
}

func (d *DockerDeployer) Stop(containerID string) error {
	logger.Deployer("Stopping container: %s", containerID)

	timeout := 10

	if err := d.client.ContainerStop(d.ctx, containerID, container.StopOptions{
		Timeout: &timeout,
	}); err != nil {
		return logger.Error("Failed to stop container %s: %w", containerID, err)
	}

	logger.Deployer("✅ Container %s stopped", containerID)

	return nil
}

func convertPorts(ports []int) (nat.PortMap, nat.PortSet) {
	bindings := nat.PortMap{}
	exposed := nat.PortSet{}

	for _, p := range ports {
		port := nat.Port(fmt.Sprintf("%d/tcp", p))
		exposed[port] = struct{}{}
		bindings[port] = []nat.PortBinding{{HostPort: fmt.Sprintf("%d", p)}}
	}

	return bindings, exposed
}

func convertVolumes(vols []string) []mount.Mount {
	mounts := make([]mount.Mount, 0, len(vols))

	for _, v := range vols {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: v,
			Target: "/data",
		})
	}

	return mounts
}

func formatEnvVars(env map[string]string) []string {
	envList := make([]string, 0, len(env))

	for key, val := range env {
		envList = append(envList, fmt.Sprintf("%s=%s", key, val))
	}

	return envList
}
