package main

import (
	"fmt"
	"os"

	"github.com/zinrai/declxc/internal/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
