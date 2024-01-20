package manager

import "fmt"

type Registry struct {
	table map[string]int
	store []Command
}

func NewRegistry(size ...uint) *Registry {
	initialSize := 0
	if len(size) > 0 {
		initialSize = int(size[0])
	}

	return &Registry{
		table: make(map[string]int),
		store: make([]Command, 0, initialSize),
	}
}

func (r *Registry) lookup(name string) bool {
	_, ok := r.table[name]
	return ok
}

func (r *Registry) get(name string) (Command, error) {
	idx, ok := r.table[name]
	if !ok {
		return nil, fmt.Errorf("command '%s' is doesn't exists", name)
	}

	return r.store[idx], nil
}

func (r *Registry) set(command Command) error {
	if _, ok := r.table[command.Name()]; ok {
		return fmt.Errorf("command '%s' already exists", command.Name())
	}

	r.table[command.Name()] = len(r.store)
	r.store = append(r.store, command)

	for _, req := range command.Requires() {
		idx, ok := r.table[req]
		if !ok {
			return fmt.Errorf("command '%s' does not exist", req)
		}

		r.store[r.table[command.Name()]].SetCommand(r.store[idx])
	}

	return nil
}
