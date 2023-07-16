package eval

import "fmt"

type Value interface {
	Get(field string) Value
}

type StructValue struct {
	// TypeParams map[Identifier]StructType
	Fields map[string]Value
}

func NewStructValue(fields map[string]Value) (sv StructValue) {
	sv.Fields = make(map[string]Value)
	for name, value := range fields {
		sv.Fields[name] = value
	}
	return
}
func NewStructValueFromType(st StructType, fields map[string]Value) (sv StructValue) {
	sv.Fields = make(map[string]Value)
	for name, _ := range st.Fields {
		var v Value
		sv.Fields[name] = v
	}
	for name, _ := range sv.Fields {
		sv.Fields[name] = fields[name]
	}
	return
}

func (sv StructValue) Get(field string) Value {
	return sv.Fields[field]
}

type Primitive struct {
	v any
	StructValue
}

func NewPrimitive(v any) Primitive {
	sv := StructValue{Fields: make(map[string]Value)}
	return Primitive{v: v, StructValue: sv}
}

var nilValue = func() (p Primitive) {
	sv := StructValue{Fields: make(map[string]Value)}
	sv.Fields["name"] = NewPrimitive("nil")
	sv.Fields["len"] = NewPrimitive(0)
	p.v = nil
	return
}()

func NewNil() (p Primitive) {
	return nilValue
}

func (p Primitive) String() string {
	return fmt.Sprintf("%v", p.v)
}

func (p Primitive) Get(field string) Value {
	if field == "v" {
		return p
	}
	switch v := p.v.(type) {
	case int:
		if field == "name" {
			return NewPrimitive("int")
		} else if field == "len" {
			return NewPrimitive(0)
		}
	case float64:
		if field == "name" {
			return NewPrimitive("float")
		} else if field == "len" {
			return NewPrimitive(0)
		}
	case bool:
		if field == "name" {
			return NewPrimitive("bool")
		} else if field == "len" {
			return NewPrimitive(0)
		}
	case string:
		if field == "name" {
			return NewPrimitive("string")
		} else if field == "len" {
			return NewPrimitive(len(v))
		}
	case nil:
		if field == "name" {
			return NewPrimitive("nil")
		} else if field == "len" {
			return NewPrimitive(0)
		}
	}
	return p.StructValue.Get(field)
}
