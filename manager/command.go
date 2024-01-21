package manager

const (
	ValidateCmdName = "validate"
	GetCmdName      = "get"
	InstallCmdName  = "install"
)

type Command interface {
	Name() string
	Requires() []string
	SetCommand(cmd Command) error
	Execute(input interface{}) (interface{}, error)

	bpmCmd()
}
