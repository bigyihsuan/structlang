package eval

import (
	"github.com/bigyihsuan/structlang/trees/ast"
)

type Evaluator struct {
	code    []ast.Stmt
	baseEnv Env
}

func NewEvaluator(code []ast.Stmt) Evaluator {
	var e Evaluator
	e.code = code
	e.baseEnv = NewEnv()
	return e
}

func (e *Evaluator) Stmt() error {
	for _, stmt := range e.code {
		switch stmt := stmt.(type) {
		case ast.TypeDef:
			return e.TypeDef(stmt)
		}
	}
	return nil
}

func (e *Evaluator) TypeDef(stmt ast.TypeDef) error {
	// TODO: implement
	return nil
}
