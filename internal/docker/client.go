package docker

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"io"
	"os"
	"os/signal"

	"golang.org/x/term"

	"Daemon/internal/shared/logger"
	"Daemon/internal/shared/utils"

	"github.com/docker/docker/api/types"
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

func (dockerClient *DockerClient) GetClient() *client.Client {
	return dockerClient.cli
}

func (dockerClient *DockerClient) CreateContainer(
	ctx context.Context,
	id, imagePath string,
	env map[string]string,
	volumes []string,
	ports []int,
	name string,
) (string, error) {
	logger.Docker("Pulling image: %s", imagePath)

	reader, err := dockerClient.cli.ImagePull(ctx, imagePath, image.PullOptions{})
	if err != nil {
		return "", logger.Error("Failed to pull image: %v", err)
	}
	io.Copy(io.Discard, reader)

	var envList []string
	for key, val := range env {
		envList = append(envList, fmt.Sprintf("%s=%s", key, val))
	}

	for _, port := range ports {
		if utils.IsPortInUse(port) {
			return "", logger.Error("Port %d is already in use on host", port)
		}
	}

	// Prepare exposed ports and bindings
	exposedPorts := nat.PortSet{}
	portBindings := nat.PortMap{}
	for _, port := range ports {
		dockerPort := nat.Port(fmt.Sprintf("%d/tcp", port))
		exposedPorts[dockerPort] = struct{}{}
		portBindings[dockerPort] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: fmt.Sprintf("%d", port),
			},
		}
	}

	containerConfig := &container.Config{
		Image:        imagePath,
		Env:          envList,
		OpenStdin:    true,
		Tty:          true,
		ExposedPorts: exposedPorts,
		Labels: map[string]string{
			"com.daemon.uuid": id,
			"com.daemon.name": name,
			"created_by":      "wings-cli",
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
	}

	resp, err := dockerClient.cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, name)
	if err != nil {
		return "", logger.Error("Failed to create container: %v", err)
	}

	logger.Docker("Container created: DockerID=%s", resp.ID)
	return resp.ID, nil
}

func (dockerClient *DockerClient) StartContainer(ctx context.Context, id string) error {
	logger.Docker("Starting container: %s", id)

	if err := dockerClient.cli.ContainerStart(ctx, id, container.StartOptions{}); err != nil {
		return logger.Error("Failed to start container %s: %v", id, err)
	}

	logger.Docker("Container started: %s", id)

	return nil
}

func (dockerClient *DockerClient) StopContainer(ctx context.Context, id string) error {
	logger.Docker("Stopping container: %s", id)

	opts := container.StopOptions{}

	if err := dockerClient.cli.ContainerStop(ctx, id, opts); err != nil {
		return logger.Error("Failed to stop container %s: %v", id, err)
	}

	logger.Docker("Container stopped: %s", id)

	return nil
}

func (dockerClient *DockerClient) RemoveContainer(ctx context.Context, id string) error {
	logger.Docker("Removing container: %s", id)

	if err := dockerClient.cli.ContainerRemove(ctx, id, container.RemoveOptions{Force: true}); err != nil {
		return logger.Error("Failed to remove container %s: %v", id, err)
	}

	logger.Docker("Container removed: %s", id)

	return nil
}

func (dockerClient *DockerClient) GetContainerByID(ctx context.Context, id string) (*container.InspectResponse, error) {
	containerJSON, err := dockerClient.cli.ContainerInspect(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("failed to inspect container %s: %w", id, err)
	}

	return &containerJSON, nil
}

func (dockerClient *DockerClient) ResolveNameToID(ctx context.Context, name string) (string, error) {
	containers, err := dockerClient.cli.ContainerList(ctx, container.ListOptions{All: true})

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

func LogOptions() container.LogsOptions {
	return container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: false,
		Follow:     true,
		Tail:       "100",
	}
}

// GetContainerLogs retrieves logs for the given container ID using the default
// log options.
func (dockerClient *DockerClient) GetContainerLogs(ctx context.Context, id string) (io.ReadCloser, error) {
	logger.Docker("Fetching logs for container: %s", id)

	reader, err := dockerClient.cli.ContainerLogs(ctx, id, LogOptions())
	if err != nil {
		return nil, logger.Error("Failed to get logs for %s: %v", id, err)
	}

	return reader, nil
}

// ExecInteractive runs a command in a container and attaches the current
// terminal so the user can interact with the process.
func (dockerClient *DockerClient) ExecInteractive(ctx context.Context, id string, cmd []string) error {
	execConfig := types.ExecConfig{
		Cmd:          cmd,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	}

	execResp, err := dockerClient.cli.ContainerExecCreate(ctx, id, execConfig)
	if err != nil {
		return logger.Error("Failed to create exec instance: %v", err)
	}

	attach, err := dockerClient.cli.ContainerExecAttach(ctx, execResp.ID, types.ExecStartCheck{Tty: true})
	if err != nil {
		return logger.Error("Failed to attach to exec instance: %v", err)
	}
	defer attach.Close()

	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return logger.Error("stdin is not a terminal — cannot enter raw mode")
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return logger.Error("Failed to set terminal raw mode: %v", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	go func() { io.Copy(attach.Conn, os.Stdin) }()
	go func() { io.Copy(os.Stdout, attach.Reader) }()

	<-sigs
	return nil
}
