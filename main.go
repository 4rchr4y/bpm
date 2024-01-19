package main

import (
	"os"

	"github.com/4rchr4y/bpm/command"
)

func main() {
	if err := command.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
