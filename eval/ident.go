package eval

import (
	"fmt"

	"github.com/bigyihsuan/structlang/trees/ast"
)

type Identifier struct {
	Name  string
	Field *Identifier
}

func NewIdentifier(name ast.Ident) Identifier {
	return Identifier{Name: name.Name, Field: nil}
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
	i.Field = &id
	return id
}

func (i Identifier) String() string {
	if i.Field != nil {
		return fmt.Sprintf("%s->%s", i.Name, i.Field.String())
	} else {
		return i.Name
	}
}
