package main

// Container represents an LXC container configuration
type Container struct {
	Name string `yaml:"name"`
	// LXCCreateArgs is passed verbatim to lxc-create (after the injected
	// "-n <name>"). Whitespace-separated; no shell quoting is interpreted.
	LXCCreateArgs string          `yaml:"lxc_create_args"`
	Networks      []NetworkConfig `yaml:"networks,omitempty"`
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
