package main

import (
	"Daemon/internal/service"
	"Daemon/pkg/logger"
	"os"
)

func main() {
	NewEgg, err := service.LoadEgg("nests/minecraft/paper.json", map[string]string{
		"MEMORY":  "1G",
		"VERSION": "latest",
		"EULA":    "TRUE",
	})

	if err != nil {
		logger.Error("Failed to load egg: %v", err)
		os.Exit(1)
	}

	deploy, err := service.NewDockerDeployer()
	if err != nil {
		logger.Error("Failed to create deployer: %v", err)
		os.Exit(1)
	}

	if err := deploy.Run(NewEgg); err != nil {
		logger.Error("Deployment failed: %v", err)
		os.Exit(1)
	}

	//time.Sleep(10 * time.Second)
	//if err := deploy.Stop("test-paper"); err != nil {
	//	logger.Error("Failed to stop container: %v", err)
	//}
}
