package validator

import (
	"fmt"
	"strings"

	"github.com/zinrai/declxc/pkg/models"
)

// ValidateContainers validates a slice of container definitions
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

			// Very basic validation for IPv4 address format if provided
			if network.IPv4Address != "" && !strings.Contains(network.IPv4Address, "/") {
				return fmt.Errorf("invalid IPv4 address format for network %d in container %s, should include CIDR notation (e.g. 192.168.2.10/24)", i, container.Name)
			}
		}

		// Validate user configuration if provided
		if len(container.Users) > 0 {
			userSet := make(map[string]bool)
			for i, user := range container.Users {
				if user.Username == "" {
					return fmt.Errorf("username is required for user %d in container %s", i, container.Name)
				}
				if user.Password == "" {
					return fmt.Errorf("password is required for user %s in container %s", user.Username, container.Name)
				}

				// Check for duplicate usernames within the same container
				if userSet[user.Username] {
					return fmt.Errorf("duplicate username %s in container %s", user.Username, container.Name)
				}
				userSet[user.Username] = true
			}
		}

		// Validate package configuration if provided
		for i, pkg := range container.Packages {
			if pkg == "" {
				return fmt.Errorf("empty package name at index %d in container %s", i, container.Name)
			}
		}
	}

	return nil
}
