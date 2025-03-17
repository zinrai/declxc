package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/zinrai/declxc/internal/lxc"
	"github.com/zinrai/declxc/internal/parser"
)

// Run initializes and runs the CLI application
func Run() error {
	// Define command-line subcommands
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createFile := createCmd.String("f", "", "Path to the YAML file with container definitions")

	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	startFile := startCmd.String("f", "", "Path to the YAML file with container definitions")

	destroyCmd := flag.NewFlagSet("destroy", flag.ExitOnError)
	destroyFile := destroyCmd.String("f", "", "Path to the YAML file with container definitions")

	// Check if a subcommand is provided
	if len(os.Args) < 2 {
		fmt.Println("Expected 'create' or 'destroy' subcommand")
		fmt.Println("Usage:")
		fmt.Println("  declxc create -f <yaml-file>")
		fmt.Println("  declxc destroy -f <yaml-file>")
		return fmt.Errorf("missing subcommand")
	}

	// Handle subcommands
	switch os.Args[1] {
	case "create":
		createCmd.Parse(os.Args[2:])
		if *createFile == "" {
			createCmd.PrintDefaults()
			return fmt.Errorf("missing file parameter")
		}
		return createContainers(*createFile)
	case "start":
		startCmd.Parse(os.Args[2:])
		if *startFile == "" {
			startCmd.PrintDefaults()
			return fmt.Errorf("missing file parameter")
		}
		return startContainers(*startFile)
	case "destroy":
		destroyCmd.Parse(os.Args[2:])
		if *destroyFile == "" {
			destroyCmd.PrintDefaults()
			return fmt.Errorf("missing file parameter")
		}
		return destroyContainers(*destroyFile)
	default:
		fmt.Printf("%q is not a valid subcommand.\n", os.Args[1])
		fmt.Println("Usage:")
		fmt.Println("  declxc create -f <yaml-file>")
		fmt.Println("  declxc destroy -f <yaml-file>")
		return fmt.Errorf("invalid subcommand")
	}
}

func createContainers(filePath string) error {
	fmt.Printf("Creating containers from %s\n", filePath)

	// Parse the YAML file
	containers, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse YAML file: %w", err)
	}

	// Create containers
	for _, container := range containers {
		if err := lxc.CreateContainer(container); err != nil {
			fmt.Printf("Error creating container %s: %v\n", container.Name, err)
		} else {
			fmt.Printf("Container %s created successfully\n", container.Name)
		}
	}

	return nil
}

func startContainers(filePath string) error {
	fmt.Printf("Starting containers from %s\n", filePath)

	// Parse the YAML file
	containers, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse YAML file: %w", err)
	}

	// Start containers
	for _, container := range containers {
		if err := lxc.StartContainer(container); err != nil {
			fmt.Printf("Error starting container %s: %v\n", container.Name, err)
		} else {
			fmt.Printf("Container %s started successfully\n", container.Name)
		}
	}

	return nil
}

func destroyContainers(filePath string) error {
	fmt.Printf("Destroying containers from %s\n", filePath)

	// Parse the YAML file
	containers, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse YAML file: %w", err)
	}

	// Destroy containers
	for _, container := range containers {
		if err := lxc.DestroyContainer(container); err != nil {
			fmt.Printf("Error destroying container %s: %v\n", container.Name, err)
		} else {
			fmt.Printf("Container %s destroyed successfully\n", container.Name)
		}
	}

	return nil
}
