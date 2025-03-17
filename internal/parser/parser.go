package parser

import (
	"fmt"
	"os"

	"github.com/zinrai/declxc/internal/validator"
	"github.com/zinrai/declxc/pkg/models"
	"gopkg.in/yaml.v3"
)

// YAML file containing container definitions
func ParseFile(filePath string) ([]models.Container, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exist", filePath)
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse YAML
	var config models.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Validate container definitions
	if err := validator.ValidateContainers(config.Containers); err != nil {
		return nil, err
	}

	return config.Containers, nil
}
