package manager

import (
	"github.com/4rchr4y/bpm/bundle"
)

type checkCmdLoader interface {
	LoadBundle(dirPath string) (*bundle.Bundle, error)
}

type (
	CheckCmdResources struct {
		FileLoader checkCmdLoader
	}

	CheckCmdInput struct {
		Path string
	}

	CheckCmdResult struct {
		Bundle *bundle.Bundle
	}
)

type checkCommand = Command[*CheckCmdResources, *CheckCmdInput, *CheckCmdResult]
type checkCommandGuardFn = func(*checkCommand, *CheckCmdInput) error

func NewCheckCommand(resources *CheckCmdResources) Commander {
	return &checkCommand{
		Name:      CheckCmdName,
		Resources: resources,
		Run:       runCheckCmd,
		Guards: []checkCommandGuardFn{
			validateCheckCmdInputPath,
		},
	}
}

func runCheckCmd(cmd *checkCommand, input *CheckCmdInput) (*CheckCmdResult, error) {
	b, err := cmd.Resources.FileLoader.LoadBundle(input.Path)
	if err != nil {
		return nil, err
	}

	if err := bundle.ValidateBundle(b); err != nil {
		return nil, err
	}

	return &CheckCmdResult{
		Bundle: b,
	}, nil
}

func validateCheckCmdInputPath(cmd *checkCommand, input *CheckCmdInput) error {
	return validatePath(input.Path)
}
