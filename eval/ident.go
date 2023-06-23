package eval

type Identifier struct {
	Name string
	Next *Identifier
}

func NewIdentifier(name string) Identifier {
	return Identifier{Name: name, Next: nil}
}

func NewIdentifierAccessing(baseName string, nextNames ...string) Identifier {
	base := NewIdentifier(baseName)
	current := base
	for _, name := range nextNames {
		next := current.NewAccess(name)
		current = next
	}
	return base
}

func (i *Identifier) NewAccess(next string) Identifier {
	id := NewIdentifier(next)
	i.Next = &id
	return id
}
