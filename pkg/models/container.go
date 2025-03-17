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
	Type        string `yaml:"type"`
	Interface   string `yaml:"interface"`
	IPv4Address string `yaml:"ipv4_address,omitempty"`
	IPv4Gateway string `yaml:"ipv4_gateway,omitempty"`
}

// Complete YAML configuration
type Config struct {
	Containers []Container `yaml:"containers"`
}
