package main

import (
	"github.com/kohkimakimoto/essh/essh"
	"os"
)

func main() {
	os.Exit(essh.Run(os.Args[1:]))
}
