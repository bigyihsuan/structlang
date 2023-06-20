package parsetree

import (
	"github.com/bigyihsuan/structlang/token"
	"github.com/bigyihsuan/structlang/util"
)

type Stmt interface{ stmtTag() }

type TypeDef struct {
	Type      token.Token
	Typename  Type
	Eq        token.Token
	StructDef Struct
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

type Struct struct {
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
