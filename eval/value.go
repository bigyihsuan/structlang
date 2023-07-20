package eval

type Value interface {
	Get(field string) Value
	TypeName() TypeName
	Unwrap() any
	PrintString() string
}
