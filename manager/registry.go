package manager

import "fmt"

type Registry struct {
	table map[string]int
	store []Command
}

func NewRegistry() *Registry {
	return &Registry{
		table: make(map[string]int),
		store: make([]Command, 0),
	}
}

func (cr *Registry) get(name string) (Command, error) {
	idx, ok := cr.table[name]
	if !ok {
		return nil, fmt.Errorf("command '%s' is doesn't exists", name)
	}

	return cr.store[idx], nil
}

func (cr *Registry) set(command Command) error {
	_, ok := cr.table[command.Name()]
	if ok {
		return fmt.Errorf("command '%s' is already exists", command.Name())
	}

	cr.table[command.Name()] = len(cr.store)
	cr.store = append(cr.store, command)

	for i := range command.Requires() {
		idx, ok := cr.table[command.Requires()[i]]
		if !ok {
			return fmt.Errorf("command '%s' is doesn't exists", command.Requires()[i])
		}

		cr.store[cr.table[command.Name()]].SetCommand(cr.store[idx])
	}

	return nil
}
