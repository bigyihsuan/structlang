package eval

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

type Primitive[T any] struct {
	v T
	StructValue
}

var zero = func() (p Primitive[int]) {
	sv := StructValue{Fields: make(map[string]Value)}
	sv.Fields["name"] = NewString("int")
	sv.Fields["len"] = zero
	p.StructValue = sv
	p.v = 0
	return
}()

func NewInt(v int) (p Primitive[int]) {
	sv := StructValue{Fields: make(map[string]Value)}
	sv.Fields["name"] = NewString("int")
	sv.Fields["len"] = zero
	p.StructValue = sv
	p.v = v
	return
}
func NewFloat(v float64) (p Primitive[float64]) {
	sv := StructValue{Fields: make(map[string]Value)}
	sv.Fields["name"] = NewString("float")
	sv.Fields["len"] = zero
	p.StructValue = sv
	p.v = v
	return
}
func NewBool(v bool) (p Primitive[bool]) {
	sv := StructValue{Fields: make(map[string]Value)}
	sv.Fields["name"] = NewString("bool")
	sv.Fields["len"] = zero
	p.StructValue = sv
	p.v = v
	return
}
func NewString(v string) (p Primitive[string]) {
	sv := StructValue{Fields: make(map[string]Value)}
	sv.Fields["name"] = NewString("string")
	sv.Fields["len"] = NewInt(len([]rune(v)))
	p.StructValue = sv
	p.v = v
	return
}

var nilValue = func() (p Primitive[*bool]) {
	sv := StructValue{Fields: make(map[string]Value)}
	sv.Fields["name"] = NewString("nil")
	sv.Fields["len"] = zero
	p.StructValue = sv
	p.v = nil
	return
}()

func NewNil() (p Primitive[*bool]) {
	return nilValue
}

func (p Primitive[T]) Get(field string) Value {
	if field == "v" {
		return p
	}
	return p.StructValue.Get(field)
}
