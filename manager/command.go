package manager

const (
	ValidateCommandName = "validate"
	GetCommandName      = "get"
)

type Command interface {
	Name() string
	Requires() []string
	SetCommand(cmd Command) error
	Execute(input interface{}) (interface{}, error)

	bpmCmd()
}
