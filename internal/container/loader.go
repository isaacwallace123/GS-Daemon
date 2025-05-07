package container

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadEgg(path string) (Egg, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Egg{}, fmt.Errorf("failed to read egg file: %w", err)
	}

	var egg Egg
	if err := json.Unmarshal(data, &egg); err != nil {
		return Egg{}, fmt.Errorf("failed to parse egg file: %w", err)
	}

	return egg, nil
}
