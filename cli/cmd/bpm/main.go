package main

import (
	"context"
	"log"
	"os"

	"github.com/4rchr4y/bpm/command/factory"
	"github.com/4rchr4y/bpm/command/root"
	"github.com/4rchr4y/bpm/internal/build"
)

type exitCode = int

const (
	exitOk exitCode = iota
	exitErr
)

func main() {
	os.Exit(run())
}

func run() exitCode {
	cmdFactory := factory.New(build.Version)

	rootCmd, err := root.NewCmdRoot(cmdFactory, build.Version)
	if err != nil {
		log.Fatalf("failed to create root command: %v\n", err)
		return exitErr
	}

	ctx := context.Background()

	if _, err := rootCmd.ExecuteContextC(ctx); err != nil {
		log.Fatal(err)

		return exitErr
	}

	return exitOk
}
