package eval

import (
	"fmt"
	"strings"
)

type Struct struct {
	Name       string
	TypeParams map[string]TypeName
	Fields     map[string]Value
}

func NewStruct(fields map[string]Value) (sv Struct) {
	sv.Fields = make(map[string]Value)
	for name, value := range fields {
		sv.Fields[name] = value
	}
	return
}
func NewStructFromType(template Type, typeParams map[string]TypeName, fields map[string]Value, typeName string) (sv Struct) {
	sv.TypeParams = typeParams
	sv.Name = typeName
	sv.Fields = make(map[string]Value)
	for name := range template.Fields {
		var v Value
		sv.Fields[name] = v
	}
	for name := range sv.Fields {
		sv.Fields[name] = fields[name]
	}
	return
}

func (sv Struct) Get(field string) Value {
	return sv.Fields[field]
}

func (sv Struct) TypeName() TypeName {
	return TypeName{Name: sv.Name}
}
func (sv Struct) Unwrap() any {
	// TODO: what is this unwrapped?
	return sv.Fields
}
func (sv Struct) printString() string {
	fields := []string{}
	for name, value := range sv.Fields {
		fields = append(fields, name+":"+value.printString())
	}
	fieldStr := ""
	if len(fields) > 0 {
		fieldStr = strings.Join(fields, ", ")
	}
	return fmt.Sprintf("%s{%s}", sv.Name, fieldStr)
}
