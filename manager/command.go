package manager

const (
	ValidateCommandName = "validate"
	GetCommandName      = "get"
	InstallCommandName  = "install"
)

type Command interface {
	Name() string
	Requires() []string
	SetCommand(cmd Command) error
	Execute(input interface{}) (interface{}, error)

	bpmCmd()
}
