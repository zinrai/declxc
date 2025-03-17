package lxc

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zinrai/declxc/pkg/models"
)

// LXC container using the provided configuration
func CreateContainer(container models.Container) error {
	// Build the lxc-create command
	args := []string{
		"-n", container.Name,
		"-t", container.Template,
		"--",
		"-r", container.Release,
		"-a", container.Arch,
	}

	cmd := exec.Command("lxc-create", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	// Configure network if specified
	if container.Network != nil {
		if err := configureNetwork(container); err != nil {
			return fmt.Errorf("failed to configure network: %w", err)
		}
	}

	return nil
}

// Destroys an LXC container
func DestroyContainer(container models.Container) error {
	cmd := exec.Command("lxc-destroy", "-n", container.Name, "-f")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to destroy container: %w", err)
	}

	return nil
}

// Container's network configuration
func configureNetwork(container models.Container) error {
	configPath := filepath.Join("/var/lib/lxc", container.Name, "config")

	// Read the current config file
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read container config: %w", err)
	}

	// Prepare network configuration
	networkConfig := fmt.Sprintf(`
# Network configuration added by declxc
lxc.net.0.type = %s
lxc.net.0.link = %s
lxc.net.0.flags = up
`, container.Network.Type, container.Network.Interface)

	// Check if network config already exists
	if strings.Contains(string(content), "lxc.net.0.type") {
		// Replace existing network config
		lines := strings.Split(string(content), "\n")
		var newLines []string
		skipNetwork := false

		for _, line := range lines {
			if strings.HasPrefix(line, "lxc.net.0.") {
				skipNetwork = true
				continue
			} else if skipNetwork && len(line) == 0 {
				skipNetwork = false
				continue
			}

			if !skipNetwork {
				newLines = append(newLines, line)
			}
		}

		// Append our network config
		newContent := strings.Join(newLines, "\n") + networkConfig
		if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
			return fmt.Errorf("failed to update container config: %w", err)
		}
	} else {
		// Append network config to the file
		f, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open container config: %w", err)
		}
		defer f.Close()

		if _, err := f.WriteString(networkConfig); err != nil {
			return fmt.Errorf("failed to write network config: %w", err)
		}
	}

	return nil
}
