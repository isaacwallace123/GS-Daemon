package service

import (
	"Daemon/pkg/logger"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"Daemon/internal/core"
)

type JSONEgg struct {
	Name    string            `json:"name"`
	Startup string            `json:"startup"`
	Env     map[string]string `json:"env"`
	Image   string            `json:"image"`
	Ports   []int             `json:"ports"`
	Volumes []string          `json:"volumes"`
}

func (e *JSONEgg) GetName() string                   { return e.Name }
func (e *JSONEgg) GetStartupCommand() string         { return e.Startup }
func (e *JSONEgg) GetEnvironment() map[string]string { return e.Env }
func (e *JSONEgg) GetImage() string                  { return e.Image }
func (e *JSONEgg) GetPorts() []int                   { return e.Ports }
func (e *JSONEgg) GetVolumes() []string              { return e.Volumes }

func LoadEgg(path string, overrides map[string]string) (core.Egg, error) {
	logger.System("Loading egg from file: %s", path)

	data, err := os.ReadFile(path)

	if err != nil {
		return nil, logger.Error("failed to read egg file: %w", err)
	}

	logger.Debug("Egg file read successfully (%d bytes)", len(data))

	var egg JSONEgg

	if err := json.Unmarshal(data, &egg); err != nil {
		return nil, logger.Error("failed to parse egg JSON: %w", err)
	}

	logger.Info("Parsed egg: %s", egg.Name)

	for k, v := range overrides {
		logger.Debug("Overriding env %s = %s", k, v)
		egg.Env[k] = v
	}

	for key, val := range egg.Env {
		placeholder := fmt.Sprintf("{{%s}}", key)

		if strings.Contains(egg.Startup, placeholder) {
			logger.Debug("Replacing %s with %s in startup command", placeholder, val)

			egg.Startup = strings.ReplaceAll(egg.Startup, placeholder, val)
		}
	}

	logger.System("Finished processing egg: %s", egg.Name)

	logger.Info(`
		🚀 Startup Summary:
		  🧩 Startup Command : %s
		  🐳 Docker Image    : %s
		  🔌 Ports           : %v
		  🗃️ Volumes         : %v
	`, egg.GetStartupCommand(), egg.GetImage(), egg.GetPorts(), egg.GetVolumes())

	return &egg, nil
}
