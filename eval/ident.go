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
func NewIdentifierFromString(name string) Identifier {
	return Identifier{Name: name, Field: nil}
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
