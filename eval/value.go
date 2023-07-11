package eval

type Value interface {
	Get(field Identifier) Value
}

type StructValue struct {
	// TypeParams map[Identifier]StructType
	Fields map[Identifier]Value
}

func (sv StructValue) Get(field Identifier) Value {
	return sv.Fields[field]
}

type Primitive[T any] struct {
	v T
}

func (iv Primitive[T]) Get(_ Identifier) Value {
	return iv
}

func NewValue[T any](val T) Primitive[T] {
	return Primitive[T]{val}
}
