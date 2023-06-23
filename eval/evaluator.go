package eval

import (
	"github.com/bigyihsuan/structlang/trees/ast"
	"github.com/kr/pretty"
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
	typename, _ := e.TypeName(currEnv, stmt.Type)
	structdef, _ := e.StructDef(currEnv, stmt.StructDef)

	currEnv.DefineType(typename.Name, structdef)
	return nil
}

func (e *Evaluator) Type(currEnv *Env, type_ ast.Type) error {

	return nil
}

func (e *Evaluator) StructDef(currEnv *Env, structDef ast.StructDef) (Struct, error) {
	var structType Struct
	structType.TypeParams = make(map[Identifier]TypeName)
	structType.Fields = make(map[Identifier]TypeName)

	for _, typeVar := range structDef.Vars.Types {
		tv, _ := e.TypeName(currEnv, typeVar)
		structType.TypeParams[tv.Name] = tv
	}

	for _, structField := range structDef.Fields {
		fieldType, _ := e.TypeName(currEnv, structField.Type)
		for _, fieldName := range structField.Names {
			name := NewIdentifier(fieldName)
			structType.Fields[name] = fieldType
		}
	}
	pretty.Println(structType)

	return structType, nil
}

func (e *Evaluator) TypeName(currEnv *Env, typename ast.Type) (TypeName, error) {
	var type_ TypeName
	name := NewIdentifier(typename.Name)
	vars := []TypeName{}
	for _, typeArg := range typename.Vars.Types {
		arg, _ := e.TypeName(currEnv, typeArg)
		vars = append(vars, arg)
	}
	type_.Name = name
	type_.Vars = vars
	return type_, nil
}
