package builtin

import (
	"fmt"
	"strings"

	. "github.com/bigyihsuan/structlang/value"
)

type Primitive struct {
	v any
	Struct
}

func NewPrimitive(v any) Primitive {
	sv := Struct{Fields: make(map[string]Value)}
	return Primitive{v: v, Struct: sv}
}

var nilValue = func() (p Primitive) {
	sv := Struct{Fields: make(map[string]Value)}
	sv.Fields["name"] = NewPrimitive("nil")
	sv.Fields["len"] = NewPrimitive(0)
	p.v = nil
	return
}()

func (p Primitive) PrintString() string {
	return fmt.Sprintf("%v", p.v)
}

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
	return p.Struct.Get(field)
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
	return IntValue{NewPrimitive(v)}
}

func (v IntValue) Get(field string) Value {
	switch field {
	case "v":
		return v
	case "name":
		return NewString("int")
	case "len":
		return NewInt(0)
	default:
		return v.Struct.Get(field)
	}
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

func (iv IntValue) Gt(other Cmp) Value {
	return NewBool(iv.Unwrap().(int) > other.Unwrap().(int))
}
func (iv IntValue) Lt(other Cmp) Value {
	return NewBool(iv.Unwrap().(int) < other.Unwrap().(int))
}
func (iv IntValue) GtEq(other Cmp) Value {
	return NewBool(iv.Unwrap().(int) >= other.Unwrap().(int))
}
func (iv IntValue) LtEq(other Cmp) Value {
	return NewBool(iv.Unwrap().(int) <= other.Unwrap().(int))
}
func (iv IntValue) Eq(other Cmp) Value {
	return NewBool(iv.Unwrap().(int) == other.Unwrap().(int))
}

type FloatValue struct {
	Primitive
}

func NewFloat(v float64) FloatValue {
	return FloatValue{NewPrimitive(v)}
}

func (v FloatValue) Get(field string) Value {
	switch field {
	case "v":
		return v
	case "name":
		return NewString("float")
	case "len":
		return NewInt(0)
	default:
		return v.Struct.Get(field)
	}
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

func (fv FloatValue) Gt(other Cmp) Value {
	return NewBool(fv.Unwrap().(float64) > other.Unwrap().(float64))
}
func (fv FloatValue) Lt(other Cmp) Value {
	return NewBool(fv.Unwrap().(float64) < other.Unwrap().(float64))
}
func (fv FloatValue) GtEq(other Cmp) Value {
	return NewBool(fv.Unwrap().(float64) >= other.Unwrap().(float64))
}
func (fv FloatValue) LtEq(other Cmp) Value {
	return NewBool(fv.Unwrap().(float64) <= other.Unwrap().(float64))
}
func (fv FloatValue) Eq(other Cmp) Value {
	return NewBool(fv.Unwrap().(float64) == other.Unwrap().(float64))
}

type BoolValue struct {
	Primitive
}

func NewBool(v bool) BoolValue {
	return BoolValue{NewPrimitive(v)}
}

func (v BoolValue) Get(field string) Value {
	switch field {
	case "v":
		return v
	case "name":
		return NewString("bool")
	case "len":
		return NewInt(0)
	default:
		return v.Struct.Get(field)
	}
}

func (bv BoolValue) Not() Value {
	return NewBool(!bv.Unwrap().(bool))
}
func (bv BoolValue) And(other Log) Value {
	return NewBool(bv.Unwrap().(bool) && other.Unwrap().(bool))
}
func (bv BoolValue) Or(other Log) Value {
	return NewBool(bv.Unwrap().(bool) || other.Unwrap().(bool))
}

type StringValue struct {
	Primitive
}

func NewString(v string) StringValue {
	return StringValue{NewPrimitive(v)}
}

func (v StringValue) Get(field string) Value {
	switch field {
	case "v":
		return v
	case "name":
		return NewString("string")
	case "len":
		return NewInt(len([]rune(v.Unwrap().(string))))
	default:
		return v.Struct.Get(field)
	}
}

func (sv StringValue) Add(other Sum) Value {
	return NewString(sv.Unwrap().(string) + other.Unwrap().(string))
}
func (sv StringValue) Sub(other Sum) Value {
	return NewString(strings.ReplaceAll(sv.Unwrap().(string), other.Unwrap().(string), ""))
}

func (sv StringValue) Gt(other Cmp) Value {
	return NewBool(sv.Unwrap().(string) > other.Unwrap().(string))
}
func (sv StringValue) Lt(other Cmp) Value {
	return NewBool(sv.Unwrap().(string) < other.Unwrap().(string))
}
func (sv StringValue) GtEq(other Cmp) Value {
	return NewBool(sv.Unwrap().(string) >= other.Unwrap().(string))
}
func (sv StringValue) LtEq(other Cmp) Value {
	return NewBool(sv.Unwrap().(string) <= other.Unwrap().(string))
}
func (sv StringValue) Eq(other Cmp) Value {
	return NewBool(sv.Unwrap().(string) == other.Unwrap().(string))
}
