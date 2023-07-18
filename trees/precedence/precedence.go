package precedence

type Precedence int

const (
	BOTTOM Precedence = iota
	SUM
	PRODUCT
	PREFIX
)
