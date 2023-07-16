package parsetree

import (
	"github.com/bigyihsuan/structlang/token"
	"github.com/bigyihsuan/structlang/util"
)

type Stmt interface{ stmtTag() }
type Expr interface{ exprTag() }

type VarDef struct {
	LetKw  token.Token
	Lvalue Lvalue
	Eq     token.Token
	Rvalue Expr
	Sc     token.Token
}

func (vd VarDef) stmtTag() {}

type VarSet struct {
	SetKw  token.Token
	Lvalue Lvalue
	Eq     token.Token
	Rvalue Expr
	Sc     token.Token
}

func (vs VarSet) stmtTag() {}

type Lvalue interface {
	Expr
	lvalueTag()
}

type FieldAccess struct {
	Lvalue Lvalue
	Arrow  token.Token
	Field  Ident
}

func (fa FieldAccess) exprTag()   {}
func (fa FieldAccess) lvalueTag() {}

type StructLiteral struct {
	TypeName Type
	Lbrace   token.Token
	Fields   SeparatedList[StructLiteralField, token.Token]
	Rbrace   token.Token
}

func (sl StructLiteral) exprTag()   {}
func (sl StructLiteral) lvalueTag() {}

type StructLiteralField struct {
	FieldName Ident
	Colon     token.Token
	Value     Expr
}

type Literal struct {
	token.Token
}

func (l Literal) exprTag()   {}
func (l Literal) lvalueTag() {}

type Ident struct {
	Name token.Token
}

func (i Ident) exprTag()   {}
func (i Ident) lvalueTag() {}

type TypeDef struct {
	TypeKw    token.Token
	TypeName  Type
	Eq        token.Token
	StructDef StructDef
	Sc        token.Token
}

func (td TypeDef) stmtTag() {}

type Type struct {
	TypeName Ident
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
	Names SeparatedList[Ident, token.Token] // ident and comma
	Type  Type
	Sc    *token.Token
}

type SeparatedList[T any, U any] []util.Pair[T, *U] // U-separated list of Ts, with optional trailing U
