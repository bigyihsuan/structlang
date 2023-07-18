package precedence

type Precedence int

const (
	BOTTOM Precedence = iota
	COMPARISON
	SUM
	PRODUCT
	PREFIX
)
