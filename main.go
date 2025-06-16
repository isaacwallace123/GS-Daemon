package main

import (
	"context"
	"log"
	"net/http"

	"Daemon/cmd/cli"
	"Daemon/internal/api/handlers"
	"Daemon/internal/api/routes"
	"Daemon/internal/container"
	"Daemon/internal/docker"
	"Daemon/internal/shared/logger"

	_ "Daemon/cmd/cli/commands"
)

func main() {
	logger.System("🟢 Starting container CLI + API...")
	ctx := context.Background()

	client, err := docker.NewDockerClient()

	if err != nil {
		logger.Error("Failed to initialize Docker client: %v", err)

		return
	}

	service := container.NewService(client)

	go func() {
		handler := &handlers.ContainerHandler{Service: service}
		mux := http.NewServeMux()
		routes.RegisterRoutes(mux, handler)
		logger.System("🚀 API Server listening on :8080")
		log.Fatal(http.ListenAndServe(":8080", mux))
	}()

	cli.RunShell(ctx, service)
}
