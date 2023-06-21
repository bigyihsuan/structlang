package ast

type Stmt interface{ stmtTag() }
type Expr interface{ exprTag() }

type TypeDef struct {
	TypeName  Type
	StructDef StructDef
}

func (td TypeDef) stmtTag() {}

type Type struct {
	TypeName string
	TypeVars TypeVars
}

type TypeVars []Type

type StructDef struct {
	TypeVars TypeVars
	Fields   []StructField
}

type StructField struct {
	Names []string
	Type  Type
}
