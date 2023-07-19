package eval

import "fmt"

type Env struct {
	Parent    *Env
	Types     map[string]Type
	Variables map[string]Value
}

func NewEnv() Env {
	return Env{
		Parent:    nil,
		Types:     make(map[string]Type),
		Variables: make(map[string]Value),
	}
}

func (e Env) MakeChild() Env {
	return Env{
		Parent:    &e,
		Types:     make(map[string]Type),
		Variables: make(map[string]Value),
	}
}

func (e *Env) DefineType(typeName string, structType Type) {
	e.Types[typeName] = structType
}
func (e Env) GetType(typeName string) *Type {
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
func (e *Env) SetVariable(name string, value Value) error {
	if variable, ok := e.Variables[name]; !ok {
		return fmt.Errorf("variable not defined: `%s`", name)
	} else if variable.TypeName().Name != value.TypeName().Name {
		return fmt.Errorf("mismatched types: want to set `%s`, got `%s`", variable.TypeName().Name, value.TypeName().Name)
	}
	e.Variables[name] = value
	return nil
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
