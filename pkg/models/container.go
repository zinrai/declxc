package models

// User represents a user account configuration for an LXC container
type User struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Shell    string `yaml:"shell,omitempty"`
}

// Container represents an LXC container configuration
type Container struct {
	Name     string          `yaml:"name"`
	Template string          `yaml:"template"`
	Release  string          `yaml:"release"`
	Arch     string          `yaml:"arch"`
	Networks []NetworkConfig `yaml:"networks,omitempty"`
	Users    []User          `yaml:"users,omitempty"`
}

// NetworkConfig represents network configuration for an LXC container
type NetworkConfig struct {
	Type        string `yaml:"type"`
	Interface   string `yaml:"interface"`
	IPv4Address string `yaml:"ipv4_address,omitempty"`
	IPv4Gateway string `yaml:"ipv4_gateway,omitempty"`
}

// Config represents the complete YAML configuration
type Config struct {
	Containers []Container `yaml:"containers"`
}
