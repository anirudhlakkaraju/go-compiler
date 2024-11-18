package object

// NewEnvironment creates a new environment
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

// NewEnclosedEnvironment creates a new environment that extends a given outer environment
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Environment holds variable bindings in the current and outer scopes
type Environment struct {
	store map[string]Object
	outer *Environment
}

// Get returns the object bindings from current or outer scope
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set stores a binding
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
