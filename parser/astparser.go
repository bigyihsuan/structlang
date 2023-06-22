package parser

import (
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
		firsttoken := stmt.TypeKw
		return ast.TypeDef{TypeName: typename, StructDef: structdef, FirstToken: firsttoken}
	}

	return
}

func (a AstParser) Type(type_ parsetree.Type) (t ast.Type) {
	typename := type_.TypeName.Lexeme()
	typevars := a.TypeVars(type_.TypeVars)
	firsttoken := type_.TypeName
	return ast.Type{TypeName: typename, TypeVars: typevars, FirstToken: firsttoken}
}

func (a AstParser) TypeVars(typeVars *parsetree.TypeVars) (tv ast.TypeVars) {
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
	var tv ast.TypeVars
	if structdef.TypeVars != nil {
		tv = a.TypeVars(structdef.TypeVars)
	}
	fields := a.StructFields(structdef.Fields)
	firsttoken := structdef.StructKw
	return ast.StructDef{TypeVars: tv, Fields: fields, FirstToken: firsttoken}
}

func (a AstParser) StructFields(fields []parsetree.StructField) (f []ast.StructField) {
	for _, field := range fields {
		f = append(f, a.StructField(field))
	}
	return f
}
func (a AstParser) StructField(field parsetree.StructField) (f ast.StructField) {
	for _, name := range field.Names {
		f.Names = append(f.Names, name.First.Lexeme())
	}
	f.Type = a.Type(field.Type)
	f.FirstToken = field.Names[0].First
	return f
}
