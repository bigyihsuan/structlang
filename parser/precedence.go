package parser

import "structlang/token"

type Precedence int

const (
	_ Precedence = iota
	LOWEST
	EQUAL
	COMPARE
	SUM
	PRODUCT
	PREFIX
	CALL
	ACCESS
)

var precedences = map[token.TokenType]Precedence{
	token.PLUS:   SUM,
	token.MINUS:  SUM,
	token.STAR:   PRODUCT,
	token.SLASH:  PRODUCT,
	token.GT:     COMPARE,
	token.LT:     COMPARE,
	token.DEQ:    EQUAL,
	token.NEQ:    EQUAL,
	token.LPAREN: CALL,
}

func (p *Parser) peekPrecedence() Precedence {
	if prec, ok := precedences[p.peek.Type]; ok {
		return prec
	} else {
		return LOWEST
	}
}
func (p *Parser) currentPrecedence() Precedence {
	if prec, ok := precedences[p.curr.Type]; ok {
		return prec
	} else {
		return LOWEST
	}
}
