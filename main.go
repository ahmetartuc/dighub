package main

import (
	"fmt"
	"os"

	"github.com/ahmetartuc/dighub/internal/cmd"
)

var version = "2.0.0"

func main() {
	if err := cmd.Execute(version); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
