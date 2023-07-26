package parsetree

import (
	"fmt"
	"strings"

	"github.com/bigyihsuan/structlang/token"
	"github.com/bigyihsuan/structlang/util"
)

type SeparatedList[T any, U any] []util.Pair[T, *U] // U-separated list of Ts, with optional trailing U

type Stmt interface {
	fmt.Stringer
	stmtTag()
}
type Expr interface {
	fmt.Stringer
	exprTag()
}

type ExprStmt struct {
	Expr Expr
	Sc   token.Token
}

func (es ExprStmt) stmtTag()       {}
func (es ExprStmt) String() string { return fmt.Sprintf("((%s) ;)", es.Expr) }

type VarDef struct {
	LetKw  token.Token
	Lvalue Lvalue
	Eq     token.Token
	Rvalue Expr
	Sc     token.Token
}

func (vd VarDef) stmtTag()       {}
func (vd VarDef) String() string { return fmt.Sprintf("(let %s = %s ;)", vd.Lvalue, vd.Rvalue) }

type VarSet struct {
	SetKw  token.Token
	Lvalue Lvalue
	Eq     token.Token
	Rvalue Expr
	Sc     token.Token
}

func (vs VarSet) stmtTag()       {}
func (vs VarSet) String() string { return fmt.Sprintf("(set %s = %s ;)", vs.Lvalue, vs.Rvalue) }

type Lvalue interface {
	Expr
	lvalueTag()
}

type FieldAccess struct {
	Lvalue Lvalue
	Arrow  token.Token
	Field  Ident
}

func (fa FieldAccess) exprTag()       {}
func (fa FieldAccess) lvalueTag()     {}
func (fa FieldAccess) String() string { return fmt.Sprintf("(-> %s %s)", fa.Lvalue, fa.Field) }

type StructLiteral struct {
	TypeName Type
	Lbrace   token.Token
	Fields   SeparatedList[StructLiteralField, token.Token]
	Rbrace   token.Token
}

func (sl StructLiteral) exprTag()   {}
func (sl StructLiteral) lvalueTag() {}
func (sl StructLiteral) String() string {
	fields := []string{}
	for _, pair := range sl.Fields {
		field := pair.First
		fields = append(fields, field.String())
	}
	return fmt.Sprintf("(%s {%s})", sl.TypeName, strings.Join(fields, " "))
}

type StructLiteralField struct {
	FieldName Ident
	Colon     token.Token
	Value     Expr
}

func (slf StructLiteralField) String() string {
	return fmt.Sprintf("(%s:%s)", slf.FieldName, slf.Value)
}

type Literal struct {
	token.Token
}

func (l Literal) exprTag()       {}
func (l Literal) lvalueTag()     {}
func (l Literal) String() string { return l.Lexeme() }

type Ident struct {
	Name token.Token
}

func (i Ident) exprTag()       {}
func (i Ident) lvalueTag()     {}
func (i Ident) String() string { return i.Name.Lexeme() }

type TypeDef struct {
	TypeKw    token.Token
	TypeName  Type
	Eq        token.Token
	StructDef StructDef
	Sc        token.Token
}

func (td TypeDef) stmtTag() {}
func (td TypeDef) String() string {
	return fmt.Sprintf("(type %s = %s ;)", td.TypeName, td.StructDef)
}

type Type struct {
	TypeName Ident
	TypeVars *TypeVars
}

func (t Type) String() string {
	return fmt.Sprintf("(%s%s)", t.TypeName, t.TypeVars)
}

type TypeVars struct {
	Lbracket token.Token
	TypeVars SeparatedList[Type, token.Token] // Type and comma
	Rbracket token.Token
}

func (tv TypeVars) String() string {
	types := []string{}
	for _, pair := range tv.TypeVars {
		ty := pair.First
		types = append(types, ty.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(types, ","))
}

type StructDef struct {
	StructKw token.Token
	TypeVars *TypeVars
	Lbrace   token.Token
	Fields   []StructField
	Rbrace   token.Token
}

func (sd StructDef) String() string {
	fields := []string{}
	for _, field := range sd.Fields {
		fields = append(fields, field.String())
	}
	return fmt.Sprintf("(struct %s {%s})", sd.TypeVars, strings.Join(fields, " "))
}

type StructField struct {
	Names SeparatedList[Ident, token.Token] // ident and comma
	Type  Type
	Sc    *token.Token
}

func (sf StructField) String() string {
	names := []string{}
	for _, pair := range sf.Names {
		name := pair.First
		names = append(names, name.String())
	}
	return fmt.Sprintf("(%s %s)", strings.Join(names, " "), sf.Type)
}

type PrefixExpr struct {
	Op    token.Token
	Right Expr
}

func (pe PrefixExpr) exprTag() {}
func (pe PrefixExpr) String() string {
	return fmt.Sprintf("(%s %s)", pe.Op.Lexeme(), pe.Right)
}

type InfixExpr struct {
	Left  Expr
	Op    token.Token
	Right Expr
}

func (ie InfixExpr) exprTag() {}
func (ie InfixExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Op.Lexeme(), ie.Left, ie.Right)
}

type GroupingExpr struct {
	Lparen token.Token
	Expr   Expr
	Rparen token.Token
}

func (ge GroupingExpr) exprTag() {}
func (ge GroupingExpr) String() string {
	return fmt.Sprintf("(%s)", ge.Expr)
}

type FuncCallExpr struct {
	Name   Lvalue
	Lparen token.Token
	Args   SeparatedList[Expr, token.Token]
	Rparen token.Token
}

func (fce FuncCallExpr) exprTag() {}
func (fce FuncCallExpr) String() string {
	args := []string{}
	for _, pair := range fce.Args {
		arg := pair.First
		args = append(args, arg.String())
	}
	return fmt.Sprintf("(%s (%s))", fce.Name, strings.Join(args, " "))
}

type FuncDef struct {
	FuncKw     token.Token
	Lparen     token.Token
	Args       SeparatedList[FuncArg, token.Token]
	Rparen     token.Token
	ReturnType *Type
	Lbrace     token.Token
	Body       []Stmt
	Rbrace     token.Token
}

func (fd FuncDef) exprTag() {}
func (fd FuncDef) String() string {
	args := []string{}
	for _, pair := range fd.Args {
		arg := pair.First
		args = append(args, arg.String())
	}
	body := []string{}
	for _, stmt := range fd.Body {
		body = append(body, stmt.String())
	}
	return fmt.Sprintf("(func (%s) %s {%s})", args, fd.ReturnType, body)
}

type FuncArg struct {
	Name Ident
	Type Type
}

func (fa FuncArg) String() string {
	return fmt.Sprintf("(%s %s)", fa.Name, fa.Type)
}

type ReturnStmt struct {
	ReturnKw token.Token
	Expr     Expr
	Sc       token.Token
}

func (rs ReturnStmt) stmtTag() {}
func (rs ReturnStmt) String() string {
	if rs.Expr != nil {
		return fmt.Sprintf("(return %s ;)", rs.Expr)
	} else {
		return "(return ;)"
	}
}
