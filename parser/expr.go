package parser

import (
	"structlang/ast"
	"structlang/token"
)

func (p *Parser) expr(precedence Precedence) ast.Expr {
	prefix := p.prefixParseFns[p.curr.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curr.Type)
		return nil
	}
	left := prefix()
	for !p.peekedTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peek.Type]
		// next token is not an infix op
		if infix == nil {
			return left
		}
		p.getNextToken()
		left = infix(left)
	}
	return left
}

func (p *Parser) prefixExpr() ast.Expr {
	expr := &ast.PrefixExpr{FirstToken: p.curr, Operator: p.curr.Lexeme}
	p.getNextToken()
	expr.Right = p.expr(PREFIX)
	expr.LastToken = p.curr
	return expr
}

func (p *Parser) infixExpr(left ast.Expr) ast.Expr {
	expr := &ast.InfixExpr{FirstToken: p.curr, Left: left, Operator: p.curr.Lexeme}
	precedence := p.currentPrecedence()
	p.getNextToken()
	expr.Right = p.expr(precedence)
	expr.LastToken = p.curr
	return expr
}
