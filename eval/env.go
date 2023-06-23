package eval

type Env struct {
	Parent    *Env
	Types     map[Identifier]Type
	Variables map[Identifier]Value
}

func NewEnv() Env {
	return Env{
		Parent:    nil,
		Types:     make(map[Identifier]Type),
		Variables: make(map[Identifier]Value),
	}
}

func (e Env) MakeChild() Env {
	return Env{
		Parent:    &e,
		Types:     make(map[Identifier]Type),
		Variables: make(map[Identifier]Value),
	}
}

func (e *Env) DefineType(name Identifier, structType Type) {
	e.Types[name] = structType
}
func (e Env) GetType(name Identifier) *Type {
	if t, ok := e.Types[name]; ok {
		return &t
	} else if e.Parent != nil {
		return e.Parent.GetType(name)
	} else {
		return nil
	}
}

func (e *Env) DefineVariable(name Identifier, value Value) {
	e.Variables[name] = value
}
func (e Env) GetVariable(name Identifier) *Value {
	if t, ok := e.Variables[name]; ok {
		return &t
	} else if e.Parent != nil {
		return e.Parent.GetVariable(name)
	} else {
		return nil
	}
}
