package lxc

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zinrai/declxc/pkg/models"
)

// CreateContainer creates an LXC container using the provided configuration
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

	// Install packages if defined
	if len(container.Packages) > 0 {
		if err := configurePackages(container); err != nil {
			return fmt.Errorf("failed to install packages: %w", err)
		}
	}

	// Configure users if defined
	if len(container.Users) > 0 {
		if err := configureUsers(container); err != nil {
			return fmt.Errorf("failed to configure users: %w", err)
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

// StartContainer starts an LXC container
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

// StopContainer stops an LXC container
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

// DestroyContainer destroys an LXC container
func DestroyContainer(container models.Container) error {
	cmd := exec.Command("lxc-destroy", "-n", container.Name, "-f")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to destroy container: %w", err)
	}

	return nil
}

// configureNetwork writes network configuration to a separate file and includes it in the main config
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

// configureUsers creates user accounts in the container using chroot
func configureUsers(container models.Container) error {
	rootfs := filepath.Join("/var/lib/lxc", container.Name, "rootfs")

	// Check if rootfs exists
	if _, err := os.Stat(rootfs); os.IsNotExist(err) {
		return fmt.Errorf("container rootfs does not exist: %s", rootfs)
	}

	// Check if any user requires sudo
	needsSudo := false
	for _, user := range container.Users {
		if user.Sudo {
			needsSudo = true
			break
		}
	}

	// Install sudo package if needed
	if needsSudo {
		fmt.Printf("Installing sudo package in container %s\n", container.Name)
		installCmd := exec.Command("chroot", rootfs, "env", "DEBIAN_FRONTEND=noninteractive",
			"apt-get", "install", "-y", "sudo")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr

		if err := installCmd.Run(); err != nil {
			return fmt.Errorf("failed to install sudo package: %w", err)
		}
	}

	for _, user := range container.Users {
		fmt.Printf("Creating user %s in container %s\n", user.Username, container.Name)

		// Check if user already exists
		checkCmd := exec.Command("chroot", rootfs, "id", user.Username)
		if err := checkCmd.Run(); err == nil {
			fmt.Printf("User %s already exists in container %s, skipping\n", user.Username, container.Name)
			continue
		}

		// Determine shell
		shell := user.Shell
		if shell == "" {
			shell = "/bin/bash"
		}

		// Create user with useradd
		useraddCmd := exec.Command("chroot", rootfs, "useradd", "-m", "-s", shell, user.Username)
		useraddCmd.Stdout = os.Stdout
		useraddCmd.Stderr = os.Stderr

		if err := useraddCmd.Run(); err != nil {
			fmt.Printf("Failed to create user %s in container %s: %v\n", user.Username, container.Name, err)
			continue
		}

		// Set password using chpasswd
		chpasswdCmd := exec.Command("chroot", rootfs, "chpasswd")
		chpasswdCmd.Stdin = strings.NewReader(fmt.Sprintf("%s:%s", user.Username, user.Password))
		chpasswdCmd.Stdout = os.Stdout
		chpasswdCmd.Stderr = os.Stderr

		if err := chpasswdCmd.Run(); err != nil {
			fmt.Printf("Failed to set password for user %s in container %s: %v\n", user.Username, container.Name, err)
			continue
		}

		fmt.Printf("User %s created successfully in container %s\n", user.Username, container.Name)

		// Configure sudo if requested
		if user.Sudo {
			if err := configureSudo(container.Name, user, rootfs); err != nil {
				fmt.Printf("Failed to configure sudo for user %s: %v\n", user.Username, err)
			}
		}

		// Configure SSH keys if specified
		if len(user.SSHKeyFiles) > 0 {
			if err := configureSSHKeys(container.Name, user, rootfs); err != nil {
				fmt.Printf("Failed to configure SSH keys for user %s: %v\n", user.Username, err)
			}
		}
	}

	return nil
}

// configureSudo sets up sudo privileges for a user
func configureSudo(containerName string, user models.User, rootfs string) error {
	if !user.Sudo {
		return nil
	}

	// Create sudoers file content
	sudoersContent := fmt.Sprintf("%s ALL=(ALL:ALL) NOPASSWD:ALL\n", user.Username)

	// Write to sudoers.d directory
	sudoersFile := filepath.Join("/etc/sudoers.d", user.Username)
	writeCmd := exec.Command("chroot", rootfs, "tee", sudoersFile)
	writeCmd.Stdin = strings.NewReader(sudoersContent)
	if err := writeCmd.Run(); err != nil {
		return fmt.Errorf("failed to write sudoers file: %w", err)
	}

	// Set proper permissions (440)
	chmodCmd := exec.Command("chroot", rootfs, "chmod", "440", sudoersFile)
	if err := chmodCmd.Run(); err != nil {
		return fmt.Errorf("failed to set sudoers file permissions: %w", err)
	}

	fmt.Printf("Sudo privileges granted to user %s in container %s (NOPASSWD)\n", user.Username, containerName)
	return nil
}

// configureSSHKeys sets up SSH public keys for a user
func configureSSHKeys(containerName string, user models.User, rootfs string) error {
	sshDir := fmt.Sprintf("/home/%s/.ssh", user.Username)
	authorizedKeysPath := fmt.Sprintf("%s/authorized_keys", sshDir)

	// Create .ssh directory
	mkdirCmd := exec.Command("chroot", rootfs, "mkdir", sshDir)
	if err := mkdirCmd.Run(); err != nil {
		return fmt.Errorf("failed to create .ssh directory: %w", err)
	}

	// Collect all SSH keys
	var allKeys []string
	for _, keyFile := range user.SSHKeyFiles {
		// Read SSH public key file
		keyContent, err := os.ReadFile(keyFile)
		if err != nil {
			fmt.Printf("Warning: Failed to read SSH key file %s: %v\n", keyFile, err)
			continue
		}

		// Trim whitespace and add to collection
		key := strings.TrimSpace(string(keyContent))
		if key != "" {
			allKeys = append(allKeys, key)
		}
	}

	// If no keys were successfully read, return
	if len(allKeys) == 0 {
		fmt.Printf("Warning: No SSH keys were successfully read for user %s\n", user.Username)
		return nil
	}

	// Write all keys to authorized_keys
	authorizedKeysContent := strings.Join(allKeys, "\n") + "\n"
	teeCmd := exec.Command("chroot", rootfs, "tee", authorizedKeysPath)
	teeCmd.Stdin = strings.NewReader(authorizedKeysContent)
	if err := teeCmd.Run(); err != nil {
		return fmt.Errorf("failed to write authorized_keys: %w", err)
	}

	// Set permissions for .ssh directory
	chmodDirCmd := exec.Command("chroot", rootfs, "chmod", "700", sshDir)
	if err := chmodDirCmd.Run(); err != nil {
		return fmt.Errorf("failed to set permissions for .ssh directory: %w", err)
	}

	// Set permissions for authorized_keys
	chmodKeysCmd := exec.Command("chroot", rootfs, "chmod", "600", authorizedKeysPath)
	if err := chmodKeysCmd.Run(); err != nil {
		return fmt.Errorf("failed to set permissions for authorized_keys: %w", err)
	}

	// Set ownership
	chownCmd := exec.Command("chroot", rootfs, "chown", "-R", fmt.Sprintf("%s:%s", user.Username, user.Username), sshDir)
	if err := chownCmd.Run(); err != nil {
		return fmt.Errorf("failed to set ownership for .ssh directory: %w", err)
	}

	fmt.Printf("SSH keys configured successfully for user %s in container %s\n", user.Username, containerName)
	return nil
}

// configurePackages installs Debian packages in the container using chroot
func configurePackages(container models.Container) error {
	rootfs := filepath.Join("/var/lib/lxc", container.Name, "rootfs")

	// Check if rootfs exists
	if _, err := os.Stat(rootfs); os.IsNotExist(err) {
		return fmt.Errorf("container rootfs does not exist: %s", rootfs)
	}

	fmt.Printf("Installing packages in container %s\n", container.Name)

	// Update package list
	fmt.Println("Updating package list...")
	updateCmd := exec.Command("chroot", rootfs, "apt-get", "update")
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr

	if err := updateCmd.Run(); err != nil {
		return fmt.Errorf("apt-get update failed: %w", err)
	}

	// Install packages
	fmt.Printf("Installing packages: %s\n", strings.Join(container.Packages, " "))
	installArgs := []string{rootfs, "env", "DEBIAN_FRONTEND=noninteractive", "apt-get", "install", "-y"}
	installArgs = append(installArgs, container.Packages...)

	installCmd := exec.Command("chroot", installArgs...)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("package installation failed: %w", err)
	}

	fmt.Printf("Packages installed successfully in container %s\n", container.Name)
	return nil
}
