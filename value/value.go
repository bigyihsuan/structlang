package value

type Value interface {
	Get(field string) Value
	TypeName() TypeName
	Return(isReturn bool) Value
	Unwrap() any
	PrintString() string
}
