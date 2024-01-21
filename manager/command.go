package manager

import (
	"fmt"
)

const (
	ValidateCmdName = "validate"
	GetCmdName      = "get"
	InstallCmdName  = "install"
)

type Commander interface {
	Set(cmd Commander) error
	GetName() string
	GetRequires() []string
}

type Command[Hub, Input, Result any] struct {
	Name     string
	Hub      Hub
	Run      func(*Command[Hub, Input, Result], Input) (Result, error)
	Requires []string
	Guards   []func(*Command[Hub, Input, Result], Input) error
	Registry *Registry
}

func (cmd *Command[H, I, R]) GetName() string       { return cmd.Name }
func (cmd *Command[H, I, R]) GetRequires() []string { return cmd.Requires }

func (cmd *Command[H, I, R]) Set(c Commander) error {
	name := c.GetName()

	if ok := cmd.Registry.lookup(name); ok {
		return fmt.Errorf("command '%s' is already installed in '%s'", name, cmd.Name)
	}

	if err := cmd.Registry.set(c); err != nil {
		return fmt.Errorf("failed to set command '%s' in '%s': %v", name, cmd.Name, err)
	}

	return nil
}

func Execute[Hub, Input, Result any](cmd Commander, input Input) (result Result, err error) {
	c, ok := cmd.(*Command[Hub, Input, Result])
	if !ok {
		return result, fmt.Errorf("type mismatch for command '%s'", cmd.GetName())
	}

	for _, guardFn := range c.Guards {
		if err = guardFn(c, input); err != nil {
			return result, fmt.Errorf("guard failed for command '%s': %v", cmd.GetName(), err)
		}
	}

	return c.Run(c, input)
}

func ExecuteInstallCmd(cmd Commander, input *InstallCmdInput) (*InstallCmdResult, error) {
	return Execute[*InstallCmdHub, *InstallCmdInput, *InstallCmdResult](cmd, input)
}

func ExecuteGetCmd(cmd Commander, input *GetCmdInput) (*GetCmdResult, error) {
	return Execute[*GetCmdHub, *GetCmdInput, *GetCmdResult](cmd, input)
}
