package validator

import (
	"fmt"

	"github.com/zinrai/declxc/pkg/models"
)

// Slice of container definitions
func ValidateContainers(containers []models.Container) error {
	if len(containers) == 0 {
		return fmt.Errorf("no container definitions found")
	}

	// Check for duplicate names
	nameSet := make(map[string]bool)
	for _, container := range containers {
		// Check for required fields
		if container.Name == "" {
			return fmt.Errorf("container name is required")
		}
		if container.Template == "" {
			return fmt.Errorf("container template is required")
		}
		if container.Release == "" {
			return fmt.Errorf("container release is required")
		}
		if container.Arch == "" {
			return fmt.Errorf("container architecture is required")
		}

		// Check for duplicate names
		if nameSet[container.Name] {
			return fmt.Errorf("duplicate container name: %s", container.Name)
		}
		nameSet[container.Name] = true

		// Validate network configuration if provided
		for i, network := range container.Networks {
			if network.Type == "" {
				return fmt.Errorf("network type is required for network %d in container %s", i, container.Name)
			}
			if network.Interface == "" {
				return fmt.Errorf("network interface is required for network %d in container %s", i, container.Name)
			}
		}
	}

	return nil
}
