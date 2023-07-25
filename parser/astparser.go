package parser

import (
	"fmt"

	"github.com/bigyihsuan/structlang/token"

	"github.com/bigyihsuan/structlang/trees/ast"
	"github.com/bigyihsuan/structlang/trees/parsetree"
)

type AstParser struct {
	tree []parsetree.Stmt
}

func NewAstParser(tree []parsetree.Stmt) AstParser {
	return AstParser{tree}
}

func (a AstParser) Parse() (stmts []ast.Stmt) {
	for _, stmt := range a.tree {
		stmts = append(stmts, a.Stmt(stmt))
	}
	return stmts
}

func (a AstParser) Stmt(stmt parsetree.Stmt) (s ast.Stmt) {
	switch stmt := stmt.(type) {
	case parsetree.TypeDef:
		typename := a.Type(stmt.TypeName)
		structdef := a.StructDef(stmt.StructDef)
		return ast.TypeDef{
			Type:      typename,
			StructDef: structdef,
			Tokens: ast.Tokens{
				FirstToken: &stmt.TypeKw,
				LastToken:  &stmt.Sc,
			},
		}
	case parsetree.VarDef:
		lvalue := a.Lvalue(stmt.Lvalue)
		rvalue := a.Expr(stmt.Rvalue)
		return ast.VarDef{
			Lvalue: lvalue,
			Rvalue: rvalue,
			Tokens: ast.Tokens{
				FirstToken: &stmt.LetKw,
				LastToken:  &stmt.Sc,
			},
		}
	case parsetree.VarSet:
		lvalue := a.Lvalue(stmt.Lvalue)
		rvalue := a.Expr(stmt.Rvalue)
		return ast.VarSet{
			Lvalue: lvalue,
			Rvalue: rvalue,
			Tokens: ast.Tokens{
				FirstToken: &stmt.SetKw,
				LastToken:  &stmt.Sc,
			},
		}
	case parsetree.ExprStmt:
		expr := a.Expr(stmt.Expr)
		return ast.ExprStmt{
			Expr: expr,
			Tokens: ast.Tokens{
				FirstToken: expr.FirstTok(),
				LastToken:  &stmt.Sc,
			},
		}
	}

	return
}

func (a AstParser) Lvalue(lv parsetree.Lvalue) (l ast.Lvalue) {
	switch lv := lv.(type) {
	case parsetree.Ident:
		return a.Ident(lv)
	case parsetree.FieldAccess:
		lvalue := a.Lvalue(lv.Lvalue)
		ident := a.Ident(lv.Field)
		firsttoken := lvalue.FirstTok()
		lasttoken := ident.LastToken
		return ast.FieldAccess{
			Lvalue: lvalue,
			Field:  ident,
			Tokens: ast.Tokens{
				FirstToken: firsttoken,
				LastToken:  lasttoken,
			},
		}
	default:
		fmt.Printf("ast unknown ast %T\n", lv)
	}
	return
}

func (a AstParser) Ident(lv parsetree.Ident) ast.Ident {
	firsttoken := lv.Name
	lasttoken := lv.Name
	return ast.Ident{
		Name: lv.Name.Lexeme(),
		Tokens: ast.Tokens{
			FirstToken: &firsttoken,
			LastToken:  &lasttoken,
		},
	}
}

func (a AstParser) Expr(expr parsetree.Expr) (e ast.Expr) {
	switch expr := expr.(type) {
	case parsetree.Literal:
		return ast.Literal{
			Token: expr.Token,
			Tokens: ast.Tokens{
				FirstToken: &expr.Token,
				LastToken:  &expr.Token,
			},
		}
	case parsetree.StructLiteral:
		return a.StructLiteral(expr)
	case parsetree.Ident:
		return a.Ident(expr)
	case parsetree.FieldAccess:
		return a.FieldAccess(expr)
	case parsetree.PrefixExpr:
		return a.PrefixExpr(expr)
	case parsetree.InfixExpr:
		return a.InfixExpr(expr)
	case parsetree.GroupingExpr:
		return a.GroupingExpr(expr)
	case parsetree.FuncCallExpr:
		return a.FuncCallExpr(expr)
	case parsetree.FuncDef:
		return a.FuncDef(expr)
	default:
		fmt.Printf("ast unknown expr %T\n", expr)
	}
	return
}

func (a AstParser) StructLiteral(expr parsetree.StructLiteral) (sl ast.StructLiteral) {
	typeName := a.Type(expr.TypeName)

	fields := []ast.StructLiteralField{}
	for _, f := range expr.Fields {
		field := a.StructLiteralField(f.First)
		fields = append(fields, field)
	}

	lastToken := expr.Rbrace
	return ast.StructLiteral{
		TypeName: typeName,
		Fields:   fields,
		Tokens: ast.Tokens{
			FirstToken: typeName.FirstToken,
			LastToken:  &lastToken,
		},
	}
}

func (a AstParser) StructLiteralField(field parsetree.StructLiteralField) (slf ast.StructLiteralField) {
	fieldName := a.Ident(field.FieldName)
	value := a.Expr(field.Value)
	return ast.StructLiteralField{
		Name:  fieldName,
		Value: value,
		Tokens: ast.Tokens{
			FirstToken: fieldName.FirstToken,
			LastToken:  value.LastTok(),
		},
	}
}

func (a AstParser) Type(type_ parsetree.Type) (t ast.Type) {
	typename := a.Ident(type_.TypeName)
	typevars := a.TypeVars(type_.TypeVars)
	firsttoken := type_.TypeName.Name
	var lasttoken token.Token
	if len(typevars) == 0 {
		lasttoken = firsttoken
	} else {
		lasttoken = type_.TypeVars.Rbracket
	}
	return ast.Type{
		Name: typename,
		Vars: typevars,
		Tokens: ast.Tokens{
			FirstToken: &firsttoken,
			LastToken:  &lasttoken,
		},
	}
}

func (a AstParser) TypeVars(typeVars *parsetree.TypeVars) (tv []ast.Type) {
	if typeVars == nil {
		return
	}
	for _, pair := range typeVars.TypeVars {
		typename := pair.First
		type_ := a.Type(typename)
		tv = append(tv, type_)
	}
	return tv
}

func (a AstParser) StructDef(structdef parsetree.StructDef) (sd ast.StructDef) {
	var tv []ast.Type
	if structdef.TypeVars != nil {
		tv = a.TypeVars(structdef.TypeVars)
	}
	fields := a.StructFields(structdef.Fields)
	firsttoken := structdef.StructKw
	lasttoken := structdef.Rbrace
	return ast.StructDef{
		Vars:   tv,
		Fields: fields,
		Tokens: ast.Tokens{
			FirstToken: &firsttoken,
			LastToken:  &lasttoken,
		},
	}
}

func (a AstParser) StructFields(fields []parsetree.StructField) (f []ast.StructField) {
	for _, field := range fields {
		f = append(f, a.StructField(field))
	}
	return f
}
func (a AstParser) StructField(field parsetree.StructField) (f ast.StructField) {
	for _, name := range field.Names {
		f.Names = append(f.Names, a.Ident(name.First))
	}
	f.Type = a.Type(field.Type)
	f.FirstToken = f.Names[0].FirstToken
	f.LastToken = f.Type.LastToken
	return f
}

func (a AstParser) FieldAccess(expr parsetree.FieldAccess) (fa ast.FieldAccess) {
	var lv ast.Lvalue
	switch l := expr.Lvalue.(type) {
	case parsetree.Ident:
		lv = a.Ident(l)
	case parsetree.FieldAccess:
		lv = a.FieldAccess(l)
	}
	f := a.Ident(expr.Field)
	return ast.FieldAccess{
		Lvalue: lv,
		Field:  f,
		Tokens: ast.Tokens{
			FirstToken: lv.FirstTok(),
			LastToken:  f.LastTok(),
		},
	}
}

func (a AstParser) PrefixExpr(expr parsetree.PrefixExpr) ast.Expr {
	right := a.Expr(expr.Right)
	return ast.PrefixExpr{
		Op:    expr.Op,
		Right: right,
		Tokens: ast.Tokens{
			FirstToken: &expr.Op,
			LastToken:  right.LastTok(),
		},
	}
}

func (a AstParser) InfixExpr(expr parsetree.InfixExpr) ast.Expr {
	left := a.Expr(expr.Left)
	right := a.Expr(expr.Right)
	return ast.InfixExpr{
		Left:  left,
		Op:    expr.Op,
		Right: right,
		Tokens: ast.Tokens{
			FirstToken: left.FirstTok(),
			LastToken:  right.FirstTok(),
		},
	}
}

func (a AstParser) GroupingExpr(expr parsetree.GroupingExpr) ast.Expr {
	return ast.GroupingExpr{
		Expr: a.Expr(expr.Expr),
		Tokens: ast.Tokens{
			FirstToken: &expr.Lparen,
			LastToken:  &expr.Rparen,
		},
	}
}

func (a AstParser) FuncCallExpr(expr parsetree.FuncCallExpr) ast.Expr {
	name := a.Lvalue(expr.Name)
	args := []ast.Expr{}
	for _, pair := range expr.Args {
		arg := a.Expr(pair.First)
		args = append(args, arg)
	}
	lastTok := name.LastTok()
	if len(args) > 0 {
		lastTok = args[len(args)-1].LastTok()
	}
	return ast.FuncCallExpr{
		Name: name,
		Args: args,
		Tokens: ast.Tokens{
			FirstToken: name.FirstTok(),
			LastToken:  lastTok,
		},
	}
}

func (a AstParser) FuncDef(expr parsetree.FuncDef) ast.Expr {
	args := []ast.FuncArg{}
	for _, arg := range expr.Args {
		argName := a.Ident(arg.First.Name)
		argType := a.Type(arg.First.Type)
		args = append(args, ast.FuncArg{
			Name: argName,
			Type: argType,
		})
	}

	body := []ast.Stmt{}
	for _, stmt := range expr.Body {
		body = append(body, a.Stmt(stmt))
	}

	return ast.FuncDef{
		Args: args,
		Body: body,
		Tokens: ast.Tokens{
			FirstToken: &expr.FuncKw,
			LastToken:  &expr.Rbrace,
		},
	}
}
