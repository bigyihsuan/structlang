package ast

import "github.com/bigyihsuan/structlang/token"

type Stmt interface{ stmtTag() }
type Expr interface{ exprTag() }

type Tokens struct {
	FirstToken, LastToken *token.Token
}

type TypeDef struct {
	TypeName  Type
	StructDef StructDef
	Tokens
}

func (td TypeDef) stmtTag() {}

type Type struct {
	TypeName string
	TypeVars TypeVars
	Tokens
}

type TypeVars struct {
	Types []Type
	Tokens
}

type StructDef struct {
	TypeVars TypeVars
	Fields   []StructField
	Tokens
}

type StructField struct {
	Names []string
	Type  Type
	Tokens
}
