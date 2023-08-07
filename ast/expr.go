package ast

import (
	"fmt"
	"structlang/token"
)

type Expr interface {
	exprTag()
	Node
}

type Lvalue interface {
	lvalueTag()
	Node
}

type PrefixExpr struct {
	FirstToken, LastToken token.Token
	Operator              string
	Right                 Expr
}

func (e PrefixExpr) exprTag()                                   {}
func (e PrefixExpr) BoundingTokens() (token.Token, token.Token) { return e.FirstToken, e.LastToken }
func (e PrefixExpr) String() string                             { return fmt.Sprintf("(%s %s)", e.Operator, e.Right) }

type InfixExpr struct {
	FirstToken, LastToken token.Token
	Left                  Expr
	Operator              string
	Right                 Expr
}

func (e InfixExpr) exprTag()                                   {}
func (e InfixExpr) BoundingTokens() (token.Token, token.Token) { return e.FirstToken, e.LastToken }
func (e InfixExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", e.Left, e.Operator, e.Right)
}
