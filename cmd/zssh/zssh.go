package main

import (
	"fmt"
	"github.com/kohkimakimoto/zssh/zssh"
	"os"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "[zssh error] %v", err)
			os.Exit(1)
		}
	}()

	if err := zssh.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "[zssh error] %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
