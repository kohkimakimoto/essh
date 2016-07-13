package main

import (
	"fmt"
	"github.com/kohkimakimoto/essh/essh"
	"github.com/kohkimakimoto/essh/support/color"
	"os"
)

func main() {
	if err := essh.Start(); err != nil {
		fmt.Fprintf(os.Stderr, color.FgRB("essh error: %v\n", err))
		os.Exit(1)
	}

	os.Exit(0)
}
