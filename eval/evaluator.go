package eval

import (
	"fmt"

	"github.com/bigyihsuan/structlang/trees/ast"
	"golang.org/x/exp/slices"
)

type Evaluator struct {
	Code    []ast.Stmt
	BaseEnv Env
}

func NewEvaluator(code []ast.Stmt) Evaluator {
	var e Evaluator
	e.Code = code
	e.BaseEnv = NewEnv()
	return e
}

func (e *Evaluator) Stmt(currEnv *Env) error {
	for _, stmt := range e.Code {
		switch stmt := stmt.(type) {
		case ast.TypeDef:
			return e.TypeDef(currEnv, stmt)
		}
	}
	return nil
}

func (e *Evaluator) TypeDef(currEnv *Env, stmt ast.TypeDef) error {
	name := NewIdentifier(stmt.Type.Name)
	vars := []Identifier{}
	for _, v := range stmt.Type.Vars.Types {
		vars = append(vars, NewIdentifier(v.Name))
	}
	structvars := []Identifier{}
	for _, v := range stmt.StructDef.Vars.Types {
		structvars = append(vars, NewIdentifier(v.Name))
	}
	fmt.Println(structvars)
	if !slices.Equal(vars, structvars) {
		return fmt.Errorf("mismatched type params: %v and %v", vars, structvars)
	}
	fields := []Field{}
	for _, field := range stmt.StructDef.Fields {
		for _, f := range field.Names {
			fname := NewIdentifier(f)
			ftype := NewIdentifier(field.Type.Name)
			fields = append(fields, Field{
				Name: fname,
				Type: ftype,
			})
		}
	}
	fmt.Println(fields)

	ty := Type{
		Name:   name,
		Vars:   vars,
		Fields: fields,
	}
	currEnv.DefineType(name, ty)

	return nil
}
