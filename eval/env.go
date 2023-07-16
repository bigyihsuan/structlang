package eval

type Env struct {
	Parent    *Env
	Types     map[string]StructType
	Variables map[string]Value
}

func NewEnv() Env {
	return Env{
		Parent:    nil,
		Types:     make(map[string]StructType),
		Variables: make(map[string]Value),
	}
}

func (e Env) MakeChild() Env {
	return Env{
		Parent:    &e,
		Types:     make(map[string]StructType),
		Variables: make(map[string]Value),
	}
}

func (e *Env) DefineType(typeName string, structType StructType) {
	e.Types[typeName] = structType
}
func (e Env) GetType(typeName string) *StructType {
	if t, ok := e.Types[typeName]; ok {
		return &t
	} else if e.Parent != nil {
		return e.Parent.GetType(typeName)
	} else {
		return nil
	}
}

func (e *Env) DefineVariable(name string, value Value) {
	e.Variables[name] = value
}
func (e Env) GetVariable(name string) *Value {
	if t, ok := e.Variables[name]; ok {
		return &t
	} else if e.Parent != nil {
		return e.Parent.GetVariable(name)
	} else {
		return nil
	}
}
