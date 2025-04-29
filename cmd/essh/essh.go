package main

import (
	"os"

	"github.com/sevir/essh/essh"
)

func main() {
	os.Exit(essh.Run(os.Args[1:]))
}
