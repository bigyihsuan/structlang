package eval

import "fmt"

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

func (p Primitive) TypeName() TypeName {
	switch p.v.(type) {
	case int:
		return TypeName{Name: "int", Vars: []TypeName{}}
	case float64:
		return TypeName{Name: "float", Vars: []TypeName{}}
	case bool:
		return TypeName{Name: "bool", Vars: []TypeName{}}
	case string:
		return TypeName{Name: "string", Vars: []TypeName{}}
	case nil:
		return TypeName{Name: "nil", Vars: []TypeName{}}
	}
	return nilValue.TypeName()
}

func (p Primitive) Unwrap() any {
	return p.v
}

type IntValue struct {
	Primitive
}

func NewInt(v int) IntValue {
	return IntValue{Primitive: NewPrimitive(v)}
}

func (iv IntValue) Pos() Value {
	iv.v = +iv.Unwrap().(int)
	return iv
}
func (iv IntValue) Neg() Value {
	iv.v = -iv.Unwrap().(int)
	return iv
}

func (iv IntValue) Add(other Sum) Value {
	r := iv.Unwrap().(int) + other.Unwrap().(int)
	return NewInt(r)
}
func (iv IntValue) Sub(other Sum) Value {
	r := iv.Unwrap().(int) - other.Unwrap().(int)
	return NewInt(r)
}
func (iv IntValue) Mul(other Product) Value {
	r := iv.Unwrap().(int) * other.Unwrap().(int)
	return NewInt(r)
}
func (iv IntValue) Div(other Product) Value {
	o := other.Unwrap().(int)
	if o == 0 {
		return NewInt(0)
	}
	r := iv.Unwrap().(int) / o
	return NewInt(r)
}

type FloatValue struct {
	Primitive
}

func NewFloat(v float64) FloatValue {
	return FloatValue{Primitive: NewPrimitive(v)}
}

func (iv FloatValue) Pos() Value {
	iv.v = +iv.v.(float64)
	return iv
}
func (iv FloatValue) Neg() Value {
	iv.v = -iv.v.(float64)
	return iv
}

func (fv FloatValue) Add(other Sum) Value {
	fv.v = fv.Unwrap().(float64) + other.Unwrap().(float64)
	return fv
}
func (fv FloatValue) Sub(other Sum) Value {
	fv.v = fv.Unwrap().(float64) - other.Unwrap().(float64)
	return fv
}
func (fv FloatValue) Mul(other Product) Value {
	fv.v = fv.Unwrap().(float64) * other.Unwrap().(float64)
	return fv
}
func (fv FloatValue) Div(other Product) Value {
	fv.v = fv.Unwrap().(float64) / other.Unwrap().(float64)
	return fv
}
