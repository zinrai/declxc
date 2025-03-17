package models

// LXC container configuration
type Container struct {
	Name     string          `yaml:"name"`
	Template string          `yaml:"template"`
	Release  string          `yaml:"release"`
	Arch     string          `yaml:"arch"`
	Networks []NetworkConfig `yaml:"networks,omitempty"`
}

// Network configuration for an LXC container
type NetworkConfig struct {
	Type      string `yaml:"type"`
	Interface string `yaml:"interface"`
}

// Complete YAML configuration
type Config struct {
	Containers []Container `yaml:"containers"`
}
