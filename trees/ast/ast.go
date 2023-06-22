package ast

import "github.com/bigyihsuan/structlang/token"

type Stmt interface {
	stmtTag()
	GetFirstToken() token.Token
}
type Expr interface{ exprTag() }

type TypeDef struct {
	TypeName   Type
	StructDef  StructDef
	FirstToken token.Token
}

func (td TypeDef) stmtTag()                   {}
func (td TypeDef) GetFirstToken() token.Token { return td.FirstToken }

type Type struct {
	TypeName   string
	TypeVars   TypeVars
	FirstToken token.Token
}

type TypeVars []Type

type StructDef struct {
	TypeVars   TypeVars
	Fields     []StructField
	FirstToken token.Token
}

type StructField struct {
	Names      []string
	Type       Type
	FirstToken token.Token
}
