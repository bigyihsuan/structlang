package eval

import (
	"fmt"
	"strings"
)

type Value interface {
	Get(field string) Value
	TypeName() TypeName
}

type StructValue struct {
	Name       string
	TypeParams map[string]TypeName
	Fields     map[string]Value
}

func NewStructValue(fields map[string]Value) (sv StructValue) {
	sv.Fields = make(map[string]Value)
	for name, value := range fields {
		sv.Fields[name] = value
	}
	return
}
func NewStructValueFromType(st StructType, typeParams map[string]TypeName, fields map[string]Value, typeName string) (sv StructValue) {
	sv.TypeParams = typeParams
	sv.Name = typeName
	sv.Fields = make(map[string]Value)
	for name := range st.Fields {
		var v Value
		sv.Fields[name] = v
	}
	for name := range sv.Fields {
		sv.Fields[name] = fields[name]
	}
	return
}

func (sv StructValue) Get(field string) Value {
	return sv.Fields[field]
}
func (sv StructValue) TypeName() TypeName {
	fields := []string{}
	for name, field := range sv.Fields {
		fields = append(fields, fmt.Sprintf("%s %s", name, field.TypeName().Name))
	}
	return TypeName{Name: fmt.Sprintf("%s{%s}", sv.Name, strings.Join(fields, "; "))}
}
