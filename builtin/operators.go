package builtin

import (
	"github.com/bigyihsuan/structlang/env"
	"github.com/bigyihsuan/structlang/trees/ast"
	. "github.com/bigyihsuan/structlang/value"
)

type Neg interface {
	Value
	Pos() Value
	Neg() Value
}

type Sum interface {
	Value
	Add(other Sum) Value
	Sub(other Sum) Value
}

type Product interface {
	Value
	Mul(other Product) Value
	Div(other Product) Value
}

type Cmp interface {
	Value
	Gt(other Cmp) Value
	Lt(other Cmp) Value
	Eq(other Cmp) Value
	GtEq(other Cmp) Value
	LtEq(other Cmp) Value
}

type Log interface {
	Value
	Not() Value
	And(other Log) Value
	Or(other Log) Value
}

type Call interface {
	Value
	Call(evaluator Eval, args ...Value) Value
}

type Eval interface {
	Evaluate(currEnv *env.Env, stmts ...[]ast.Stmt) error
}
