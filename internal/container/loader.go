package container

import (
	"Daemon/internal/models"
	"Daemon/internal/shared/logger"

	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func LoadEgg(path string) (models.Egg, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return models.Egg{}, fmt.Errorf("failed to read egg file: %w", err)
	}

	var egg models.Egg
	if err := json.Unmarshal(data, &egg); err != nil {
		return models.Egg{}, fmt.Errorf("failed to parse egg file: %w", err)
	}

	return egg, nil
}

func FindEggPathByName(eggName string) (string, error) {
	var result string

	search := eggName + ".egg.json"

	err := filepath.Walk("nests", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return logger.Error("failed to access path during walk: %v", err)
		}

		if !info.IsDir() && strings.EqualFold(info.Name(), search) {
			result = path

			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return "", logger.Error("failed during file walk: %v", err)
	}

	if result == "" {
		return "", logger.Error("egg not found: %s", eggName)
	}

	return result, nil
}
