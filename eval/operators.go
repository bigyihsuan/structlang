package eval

type Neg interface {
	Value
	Pos() Value
	Neg() Value
}

type Sum interface {
	Value
	Add(other Sum) Value
	Sub(other Sum) Value
}

type Product interface {
	Value
	Mul(other Product) Value
	Div(other Product) Value
}
