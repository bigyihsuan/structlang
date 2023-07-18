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

type Cmp interface {
	Value
	Gt(other Cmp) Value
	Lt(other Cmp) Value
	Eq(other Cmp) Value
	GtEq(other Cmp) Value
	LtEq(other Cmp) Value
}

type Log interface {
	Not() Value
	And(other Log) Value
	Or(other Log) Value
}
