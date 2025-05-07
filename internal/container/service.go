package container

import (
	"Daemon/internal/docker"
	"Daemon/internal/shared/logger"

	"context"

	"github.com/google/uuid"
)

type Service struct {
	client docker.Client

	store map[string]*Container
}

func NewService(client docker.Client) *Service {
	return &Service{
		client: client,
		store:  make(map[string]*Container),
	}
}

func (service *Service) CreateContainer(ctx context.Context, egg Egg) (*Container, error) {
	id := uuid.New().String()

	dockerID, err := service.client.CreateContainer(ctx, id, egg.Image, egg.Env, egg.Volumes, egg.Ports)

	if err != nil {
		return nil, err
	}

	c := &Container{
		ID:       id,
		Name:     egg.Name,
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

	service.store[id] = c

	return c, nil
}

func (service *Service) StartContainer(ctx context.Context, id string) error {
	if err := service.client.StartContainer(ctx, id); err != nil {
		return err
	}

	if container, ok := service.store[id]; ok {
		container.Status = "running"
	}

	return nil
}

func (service *Service) StopContainer(ctx context.Context, id string) error {
	if err := service.client.StopContainer(ctx, id); err != nil {
		return err
	}

	if container, ok := service.store[id]; ok {
		container.Status = "stopped"
	}

	return nil
}

func (service *Service) RemoveContainer(ctx context.Context, id string) error {
	if err := service.client.RemoveContainer(ctx, id); err != nil {
		return err
	}

	delete(service.store, id)

	return nil
}

func (service *Service) GetContainer(id string) (*Container, error) {
	container, ok := service.store[id]

	if !ok {
		return nil, logger.Error("Container not found: %s", id)
	}

	return container, nil
}
