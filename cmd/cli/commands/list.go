package commands

import (
	"fmt"
	"github.com/docker/docker/api/types/filters"
	"strings"

	"Daemon/cmd/cli"
	"Daemon/internal/docker"
	"Daemon/internal/shared/logger"
	"Daemon/internal/shared/utils"

	"github.com/docker/docker/api/types/container"
)

func init() {
	cli.RegisterCommand(&cli.Command{
		Name:        "list",
		Description: "List containers created by this system (wings-cli)",
		Args:        []cli.ArgSpec{},
		Execute:     runList,
	})
}

func runList(command *cli.CommandContext) error {
	client, err := docker.NewDockerClient()
	if err != nil {
		return logger.Error("Failed to create Docker client: %v", err)
	}

	filter := filters.NewArgs()
	filter.Add("label", "created_by=wings-cli")

	containers, err := client.GetClient().ContainerList(command.Ctx, container.ListOptions{
		All:     true,
		Filters: filter,
	})

	if err != nil {
		return logger.Error("Failed to list containers: %v", err)
	}

	if len(containers) == 0 {
		fmt.Println("No containers created by wings-cli.")
		return nil
	}

	const (
		idWidth    = 12
		nameWidth  = 20
		imageWidth = 22
		portWidth  = 22
		stateWidth = 12
	)

	fmt.Printf("\n%-*s  %-*s  %-*s  %-*s  %-*s\n",
		idWidth, "ID",
		nameWidth, "Name",
		imageWidth, "Image",
		portWidth, "Port(s)",
		stateWidth, "Status")

	fmt.Printf("%s  %s  %s  %s  %s\n",
		strings.Repeat("-", idWidth),
		strings.Repeat("-", nameWidth),
		strings.Repeat("-", imageWidth),
		strings.Repeat("-", portWidth),
		strings.Repeat("-", stateWidth))

	for _, ctr := range containers {
		ports := "-"
		if len(ctr.Ports) > 0 {
			var pairs []string
			for _, p := range ctr.Ports {
				pairs = append(pairs, fmt.Sprintf("%s:%d->%d/%s", p.IP, p.PublicPort, p.PrivatePort, p.Type))
			}
			ports = strings.Join(pairs, ",")
		}

		fmt.Printf("%-*s  %-*s  %-*s  %-*s  %-*s\n",
			idWidth, ctr.ID[:12],
			nameWidth, ctr.Labels["com.daemon.name"],
			imageWidth, utils.TruncateMiddle(ctr.Image, imageWidth),
			portWidth, utils.TruncateMiddle(ports, portWidth),
			stateWidth, ctr.State)
	}

	fmt.Println()
	return nil
}
