package precedence

type Precedence int

const (
	BOTTOM Precedence = iota
	LOGICAL
	COMPARISON
	SUM
	PRODUCT
	PREFIX
)
