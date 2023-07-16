package eval

type Env struct {
	Parent    *Env
	Types     map[Identifier]StructType
	Variables map[Identifier]Value
}

func NewEnv() Env {
	return Env{
		Parent:    nil,
		Types:     make(map[Identifier]StructType),
		Variables: make(map[Identifier]Value),
	}
}

func (e Env) MakeChild() Env {
	return Env{
		Parent:    &e,
		Types:     make(map[Identifier]StructType),
		Variables: make(map[Identifier]Value),
	}
}

func (e *Env) DefineType(typeName Identifier, structType StructType) {
	e.Types[typeName] = structType
}
func (e Env) GetType(typeName Identifier) *StructType {
	if t, ok := e.Types[typeName]; ok {
		return &t
	} else if e.Parent != nil {
		return e.Parent.GetType(typeName)
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
