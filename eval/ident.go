package eval

import "github.com/bigyihsuan/structlang/trees/ast"

type Identifier struct {
	Name string
	Next *Identifier
}

func NewIdentifier(name ast.Ident) Identifier {
	return Identifier{Name: name.Name, Next: nil}
}

func NewIdentifierAccessing(baseName ast.Ident, nextNames ...ast.Ident) Identifier {
	base := NewIdentifier(baseName)
	current := base
	for _, name := range nextNames {
		next := current.NewAccess(name)
		current = next
	}
	return base
}

func (i *Identifier) NewAccess(next ast.Ident) Identifier {
	id := NewIdentifier(next)
	i.Next = &id
	return id
}
