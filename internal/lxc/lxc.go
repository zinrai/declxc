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
	// Check if the container already exists
	exists, err := containerExists(container.Name)
	if err != nil {
		return fmt.Errorf("failed to check if container exists: %w", err)
	}

	//  Create container only if it does not exist
	if !exists {
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

		fmt.Printf("Container %s created successfully\n", container.Name)
	} else {
		fmt.Printf("Container %s already exists, skipping creation\n", container.Name)
	}

	// Configure network regardless of whether container was just created or already existed
	if len(container.Networks) > 0 {
		if err := configureNetwork(container); err != nil {
			return fmt.Errorf("failed to configure network: %w", err)
		}
	}

	return nil
}

func containerExists(name string) (bool, error) {
	cmd := exec.Command("lxc-ls", "-1")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to list containers: %w", err)
	}

	containers := strings.Split(string(output), "\n")
	for _, c := range containers {
		if c == name {
			return true, nil
		}
	}

	return false, nil
}

// Starts an LXC container
func StartContainer(container models.Container) error {
	// Check if container exists
	exists, err := containerExists(container.Name)
	if err != nil {
		return fmt.Errorf("failed to check if container exists: %w", err)
	}

	if !exists {
		return fmt.Errorf("container %s does not exist", container.Name)
	}

	// Execute lxc-start command
	cmd := exec.Command("lxc-start", "-n", container.Name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil
}

func StopContainer(container models.Container) error {
	// Check if container exists
	exists, err := containerExists(container.Name)
	if err != nil {
		return fmt.Errorf("failed to check if container exists: %w", err)
	}

	if !exists {
		return fmt.Errorf("container %s does not exist", container.Name)
	}

	// Execute lxc-stop command
	cmd := exec.Command("lxc-stop", "-n", container.Name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
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

// Writes network configuration to a separate file and includes it in the main config
func configureNetwork(container models.Container) error {
	// Network config file path
	containerPath := filepath.Join("/var/lib/lxc", container.Name)
	networkConfigPath := filepath.Join(containerPath, "config-network")
	mainConfigPath := filepath.Join(containerPath, "config")

	// Create network config content
	var networkContent strings.Builder
	networkContent.WriteString("# Network configuration added by declxc\n")

	// Support multiple network interfaces
	for i, network := range container.Networks {
		networkContent.WriteString(fmt.Sprintf(`
lxc.net.%d.type = %s
lxc.net.%d.link = %s
lxc.net.%d.flags = up
`, i, network.Type, i, network.Interface, i))

		// Add IPv4 address if provided
		if network.IPv4Address != "" {
			networkContent.WriteString(fmt.Sprintf("lxc.net.%d.ipv4.address = %s\n", i, network.IPv4Address))
		}

		// Add IPv4 gateway if provided
		if network.IPv4Gateway != "" {
			networkContent.WriteString(fmt.Sprintf("lxc.net.%d.ipv4.gateway = %s\n", i, network.IPv4Gateway))
		}
	}

	// Write network config to file
	if err := os.WriteFile(networkConfigPath, []byte(networkContent.String()), 0644); err != nil {
		return fmt.Errorf("failed to write network config: %w", err)
	}

	// Update main config file to include network config
	mainContent, err := os.ReadFile(mainConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read container config: %w", err)
	}

	// Check if include line already exists
	includeLine := fmt.Sprintf("lxc.include = %s", networkConfigPath)
	if strings.Contains(string(mainContent), includeLine) {
		// Already configured, nothing to do
		return nil
	}

	// Add include line to main config
	updatedContent := string(mainContent) + "\n" + "# Network configuration include\n" + includeLine + "\n"

	if err := os.WriteFile(mainConfigPath, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to update container config: %w", err)
	}

	return nil
}
