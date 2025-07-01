package container

import (
	"Daemon/internal/docker"
	"Daemon/internal/models"
	"Daemon/internal/shared/logger"
	"context"
	"github.com/google/uuid"
	"io"
	"os"
)

type Service struct {
	client docker.Client

	store map[string]*models.Container
}

func NewService(client docker.Client) *Service {
	return &Service{
		client: client,
		store:  make(map[string]*models.Container),
	}
}

func (service *Service) CreateContainer(ctx context.Context, name string, eggName string) (*models.Container, error) {
	eggPath, err := FindEggPathByName(eggName)

	if err != nil {
		return nil, err
	}

	egg, err := LoadEgg(eggPath)

	if err != nil {
		return nil, err
	}

	id := uuid.New().String()

	dockerID, err := service.client.CreateContainer(ctx, id, egg.Image, egg.Env, egg.Volumes, egg.Ports, name)

	if err != nil {
		return nil, err
	}

	newContainer := &models.Container{
		ID:       id,
		Name:     name,
		Image:    egg.Image,
		DockerID: dockerID,
		Status:   "created",
	}

	logger.Info(`
🚀 Startup Summary:
  🧩 Startup Command : %s
  🐳 Docker Image    : %s
  🔌 Ports           : %v
  🗃️ Volumes         : %v`, egg.Startup, egg.Image, egg.Ports, egg.Volumes)

	service.store[id] = newContainer

	return newContainer, nil
}

func (service *Service) StartContainer(ctx context.Context, name string) error {
	id, err := service.client.ResolveNameToID(ctx, name)
	if err != nil {
		return err
	}
	return service.client.StartContainer(ctx, id)
}

func (service *Service) StopContainer(ctx context.Context, name string) error {
	id, err := service.client.ResolveNameToID(ctx, name)
	if err != nil {
		return err
	}
	return service.client.StopContainer(ctx, id)
}

func (service *Service) RemoveContainer(ctx context.Context, name string) error {
	id, err := service.client.ResolveNameToID(ctx, name)
	if err != nil {
		return err
	}
	return service.client.RemoveContainer(ctx, id)
}

func (service *Service) GetContainer(ctx context.Context, name string) (*models.Container, error) {
	id, err := service.client.ResolveNameToID(ctx, name)
	if err != nil {
		return nil, err
	}

	inspect, err := service.client.GetContainerByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.Container{
		ID:       id,
		Name:     name,
		Image:    inspect.Config.Image,
		DockerID: id,
		Status:   inspect.State.Status,
	}, nil
}

// ExecCommand executes a command interactively inside the specified container.
func (service *Service) ExecCommand(ctx context.Context, name string, cmd []string) error {
	id, err := service.client.ResolveNameToID(ctx, name)
	if err != nil {
		return err
	}
	return service.client.ExecInteractive(ctx, id, cmd)
}

// StreamLogs prints the latest logs from the specified container and
// continues following them until interrupted.
func (service *Service) StreamLogs(ctx context.Context, name string) error {
	id, err := service.client.ResolveNameToID(ctx, name)
	if err != nil {
		return err
	}
	reader, err := service.client.GetContainerLogs(ctx, id)
	if err != nil {
		return err
	}
	defer reader.Close()
	_, err = io.Copy(os.Stdout, reader)
	return err
}
