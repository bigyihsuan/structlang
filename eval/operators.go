package eval

type Neg interface {
	Pos() Value
	Neg() Value
}

type Sum interface {
	Add(other Sum) Value
	Sub(other Sum) Value
}

type Product interface {
	Mul(other Product) Value
	Div(other Product) Value
}
