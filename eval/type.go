package eval

type TypeName struct {
	Name Identifier
	Vars []TypeName
}

type Struct struct {
	TypeParams map[Identifier]TypeName
	Fields     map[Identifier]TypeName
}

type StructInstance struct {
	TypeParams map[Identifier]Struct
	Fields     map[Identifier]Value
}
