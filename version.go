package main

import "fmt"

// Build information injected at release time via goreleaser ldflags
// (see .goreleaser.yaml). Defaults apply to a local `go build`.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// printVersion writes the build information to stdout.
func printVersion() {
	fmt.Printf("declxc %s\n", version)
	fmt.Printf("commit: %s\n", commit)
	fmt.Printf("built:  %s\n", date)
}
