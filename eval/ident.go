package eval

type Identifier struct {
	Name string
	Next *Identifier
}
