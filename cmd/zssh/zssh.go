package main

import (
	"log"
	"os"
	"github.com/kohkimakimoto/zssh/zssh"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Print("Error: %s\n", err)
			os.Exit(1)
		}
	}()

	os.Exit(zssh.Main())
}
