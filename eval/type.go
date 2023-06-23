package eval

type Type struct {
	Name   Identifier
	Vars   []Type
	Fields []Field
}

type Field struct {
	Name Identifier
	Type Type
}
