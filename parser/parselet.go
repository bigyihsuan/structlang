package parser

import (
	"errors"
	"fmt"

	"github.com/bigyihsuan/structlang/token"
	"github.com/bigyihsuan/structlang/trees/parsetree"
	"github.com/bigyihsuan/structlang/trees/precedence"
	"github.com/bigyihsuan/structlang/util"
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

type CallParselet struct{}

func (cp CallParselet) Parse(parser *ParseTreeParser, expr parsetree.Expr, lparen token.Token) (parsetree.Expr, error) {
	args := parsetree.SeparatedList[parsetree.Expr, token.Token]{}
	funcName, isLvalue := expr.(parsetree.Lvalue)
	if !isLvalue {
		return expr, errors.New("expected lvalue for function call")
	}

	if hasRparen, err := parser.nextTokenIs(token.RPAREN); err != nil {
		return funcName, err
	} else if !hasRparen {
		for {
			if finishFuncCall, err := parser.nextTokenIs(token.RPAREN); err != nil {
				return funcName, err
			} else if finishFuncCall {
				break
			}
			arg, err := parser.Expr(precedence.BOTTOM)
			if err != nil {
				return arg, err
			}
			if finishFuncCall, err := parser.nextTokenIs(token.RPAREN); err != nil {
				return funcName, err
			} else if finishFuncCall {
				args = append(args, util.Pair[parsetree.Expr, *token.Token]{First: arg, Last: nil})
				break
			}
			comma, err := parser.expectGet(token.COMMA)
			if err != nil {
				return arg, err
			}
			args = append(args, util.Pair[parsetree.Expr, *token.Token]{First: arg, Last: comma})
		}
	}
	rparen, err := parser.expectGet(token.RPAREN)
	if err != nil {
		return funcName, err
	}
	return parsetree.FuncCallExpr{Name: funcName, Lparen: lparen, Args: args, Rparen: *rparen}, nil
}
func (cp CallParselet) Precedence() precedence.Precedence { return precedence.CALL }

type FuncDefParselet struct{}

func (fdp FuncDefParselet) Parse(parser *ParseTreeParser, op token.Token) (parsetree.Expr, error) {
	fderr := errors.New("in funcdef")
	funcKw := op
	lparen, err := parser.expectGet(token.LPAREN)
	if err != nil {
		return nil, errors.Join(fderr, err)
	}
	args := parsetree.SeparatedList[parsetree.FuncArg, token.Token]{}

	for {
		if finishFuncArgs, err := parser.nextTokenIs(token.RPAREN); err != nil {
			return nil, errors.Join(fderr, err)
		} else if finishFuncArgs {
			break
		}
		argName, err := parser.Ident()
		if err != nil {
			return nil, errors.Join(fderr, err)
		}
		argType, err := parser.Type()
		if err != nil {
			return nil, errors.Join(fderr, err)
		}
		arg := parsetree.FuncArg{Name: argName, Type: argType}
		if finishFuncArgs, err := parser.nextTokenIs(token.RPAREN); err != nil {
			return nil, errors.Join(fderr, err)
		} else if finishFuncArgs {
			args = append(args, util.Pair[parsetree.FuncArg, *token.Token]{First: arg, Last: nil})
			break
		}
		comma, err := parser.expectGet(token.COMMA)
		if err != nil {
			return nil, errors.Join(fderr, err)
		}
		args = append(args, util.Pair[parsetree.FuncArg, *token.Token]{First: arg, Last: comma})
	}
	rparen, err := parser.expectGet(token.RPAREN)
	if err != nil {
		return nil, errors.Join(fderr, err)
	}
	var returnType parsetree.Type
	if hasReturnType, err := parser.nextTokenIsAny(token.IDENT, token.FUNC); err != nil {
		return nil, errors.Join(fderr, err)
	} else if hasReturnType {
		rt, err := parser.Type()
		if err != nil {
			return nil, errors.Join(fderr, err)
		}
		returnType = rt
	}
	lbrace, err := parser.expectGet(token.LBRACE)
	if err != nil {
		return nil, errors.Join(fderr, err)
	}
	body := []parsetree.Stmt{}
	for {
		if finishFuncBody, err := parser.nextTokenIs(token.RBRACE); err != nil {
			return nil, errors.Join(fderr, err)
		} else if finishFuncBody {
			break
		}
		stmt, err := parser.Stmt()
		if err != nil {
			return nil, errors.Join(fderr, err)
		}
		body = append(body, stmt)
	}
	rbrace, err := parser.expectGet(token.RBRACE)
	if err != nil {
		return nil, errors.Join(fderr, err)
	}
	return parsetree.FuncDef{
		FuncKw:     funcKw,
		Lparen:     *lparen,
		Args:       args,
		Rparen:     *rparen,
		ReturnType: returnType,
		Lbrace:     *lbrace,
		Body:       body,
		Rbrace:     *rbrace,
	}, err
}
