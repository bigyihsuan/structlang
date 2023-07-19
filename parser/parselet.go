package parser

import (
	"errors"
	"fmt"

	"github.com/bigyihsuan/structlang/token"
	"github.com/bigyihsuan/structlang/trees/parsetree"
	"github.com/bigyihsuan/structlang/trees/precedence"
)

type PrefixParselet interface {
	Parse(parser *ParseTreeParser, op token.Token) (parsetree.Expr, error)
}
type InfixParselet interface {
	Parse(parser *ParseTreeParser, left parsetree.Expr, op token.Token) (parsetree.Expr, error)
	Precedence() precedence.Precedence
}

type LiteralParselet struct{}

func (lp LiteralParselet) Parse(parser *ParseTreeParser, tok token.Token) (parsetree.Expr, error) {
	return parsetree.Literal{Token: tok}, nil
}

type IdentParselet struct{}

func (ip IdentParselet) Parse(parser *ParseTreeParser, tok token.Token) (parsetree.Expr, error) {
	parser.putBackToken()
	return parser.IdentOrStructLiteralOrFieldAccess()
}

type PrefixOperator struct {
	prec precedence.Precedence
}

func (pop PrefixOperator) Parse(parser *ParseTreeParser, op token.Token) (parsetree.Expr, error) {
	poperr := fmt.Errorf("in prefix operator `%s`", op.String())
	right, err := parser.Expr(pop.prec)
	if err != nil {
		return right, errors.Join(poperr, err)
	}
	return parsetree.PrefixExpr{Op: op, Right: right}, nil
}
func (pop PrefixOperator) Precedence() precedence.Precedence {
	return pop.prec
}

type GroupingParselet struct{}

func (gp GroupingParselet) Parse(parser *ParseTreeParser, lparen token.Token) (parsetree.Expr, error) {
	expr, err := parser.Expr(precedence.BOTTOM)
	if err != nil {
		return expr, err
	}
	rparen, err := parser.expectGet(token.RPAREN)
	if err != nil {
		return expr, err
	}
	return parsetree.GroupingExpr{Lparen: lparen, Expr: expr, Rparen: *rparen}, nil
}

type BinaryOperator struct {
	prec    precedence.Precedence
	isRight bool
}

func (bop BinaryOperator) Parse(parser *ParseTreeParser, left parsetree.Expr, op token.Token) (parsetree.Expr, error) {
	boperr := fmt.Errorf("in infix operator `%s`", op.String())
	prec := bop.prec
	if bop.isRight {
		prec -= 1
	}
	right, err := parser.Expr(prec)
	if err != nil {
		return right, errors.Join(boperr, err)
	}
	return parsetree.InfixExpr{Left: left, Op: op, Right: right}, err
}

func (bop BinaryOperator) Precedence() precedence.Precedence {
	return bop.prec
}
