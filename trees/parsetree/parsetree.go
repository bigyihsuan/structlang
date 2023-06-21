package parsetree

import (
	"github.com/bigyihsuan/structlang/token"
	"github.com/bigyihsuan/structlang/util"
)

type Stmt interface{ stmtTag() }
type Expr interface{ exprTag() }

type TypeDef struct {
	TypeKw    token.Token
	TypeName  Type
	Eq        token.Token
	StructDef StructDef
	Sc        token.Token
}

func (td TypeDef) stmtTag() {}

type Type struct {
	TypeName token.Token
	TypeVars *TypeVars
}

type TypeVars struct {
	Lbracket token.Token
	TypeVars SeparatedList[Type, token.Token] // Type and comma
	Rbracket token.Token
}

type StructDef struct {
	StructKw token.Token
	TypeVars *TypeVars
	Lbrace   token.Token
	Fields   []StructField
	Rbrace   token.Token
}

type StructField struct {
	Names SeparatedList[token.Token, token.Token] // ident and comma
	Type  Type
	Sc    *token.Token
}

type SeparatedList[T any, U any] []util.Pair[T, *U] // U-separated list of Ts, with optional trailing U
