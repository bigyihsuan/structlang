package ast

import "github.com/bigyihsuan/structlang/token"

type HasTokens interface {
	FirstTok() *token.Token
	LastTok() *token.Token
}
type Stmt interface {
	HasTokens
	stmtTag()
}
type Expr interface {
	HasTokens
	exprTag()
}

type Tokens struct {
	FirstToken, LastToken *token.Token
}

type TypeDef struct {
	Type      Type
	StructDef StructDef
	Tokens
}

func (td TypeDef) stmtTag()               {}
func (td TypeDef) FirstTok() *token.Token { return td.FirstToken }
func (td TypeDef) LastTok() *token.Token  { return td.LastToken }

type VarDef struct {
	Lvalue Lvalue
	Rvalue Expr
	Tokens
}

func (vd VarDef) stmtTag()               {}
func (vd VarDef) FirstTok() *token.Token { return vd.FirstToken }
func (vd VarDef) LastTok() *token.Token  { return vd.LastToken }

type VarSet struct {
	Lvalue Lvalue
	Rvalue Expr
	Tokens
}

func (vs VarSet) stmtTag()               {}
func (vs VarSet) FirstTok() *token.Token { return vs.FirstToken }
func (vs VarSet) LastTok() *token.Token  { return vs.LastToken }

type Lvalue interface {
	Expr
	lvalueTag()
}

type FieldAccess struct {
	Lvalue Lvalue
	Field  Ident
	Tokens Tokens
}

func (fa FieldAccess) exprTag()               {}
func (fa FieldAccess) lvalueTag()             {}
func (fa FieldAccess) FirstTok() *token.Token { return fa.Lvalue.FirstTok() }
func (fa FieldAccess) LastTok() *token.Token  { return fa.Field.LastToken }

type Type struct {
	Name Ident
	Vars []Type
	Tokens
}

func (t Type) FirstTok() *token.Token { return t.FirstToken }
func (t Type) LastTok() *token.Token  { return t.LastToken }

type StructDef struct {
	Vars   []Type
	Fields []StructField
	Tokens
}

func (sd StructDef) FirstTok() *token.Token { return sd.FirstToken }
func (sd StructDef) LastTok() *token.Token  { return sd.LastToken }

type StructField struct {
	Names []Ident
	Type  Type
	Tokens
}

func (sf StructField) FirstTok() *token.Token { return sf.FirstToken }
func (sf StructField) LastTok() *token.Token  { return sf.LastToken }

type Ident struct {
	Name string
	Tokens
}

func (i Ident) exprTag()               {}
func (i Ident) lvalueTag()             {}
func (i Ident) FirstTok() *token.Token { return i.FirstToken }
func (i Ident) LastTok() *token.Token  { return i.LastToken }

type StructLiteral struct {
	TypeName Type
	Fields   []StructLiteralField
	Tokens
}

func (sl StructLiteral) exprTag()               {}
func (sl StructLiteral) FirstTok() *token.Token { return sl.FirstToken }
func (sl StructLiteral) LastTok() *token.Token  { return sl.LastToken }

type StructLiteralField struct {
	Name  Ident
	Value Expr
	Tokens
}

func (slf StructLiteralField) FirstTok() *token.Token { return slf.FirstToken }
func (slf StructLiteralField) LastTok() *token.Token  { return slf.LastToken }

type Literal struct {
	token.Token
	Tokens
}

func (l Literal) exprTag()               {}
func (l Literal) FirstTok() *token.Token { return l.FirstToken }
func (l Literal) LastTok() *token.Token  { return l.LastToken }

type PrefixExpr struct {
	Op    token.Token
	Right Expr
	Tokens
}

func (pe PrefixExpr) exprTag()               {}
func (pe PrefixExpr) FirstTok() *token.Token { return pe.FirstToken }
func (pe PrefixExpr) LastTok() *token.Token  { return pe.LastToken }

type InfixExpr struct {
	Left  Expr
	Op    token.Token
	Right Expr
	Tokens
}

func (ie InfixExpr) exprTag()               {}
func (ie InfixExpr) FirstTok() *token.Token { return ie.FirstToken }
func (ie InfixExpr) LastTok() *token.Token  { return ie.LastToken }
