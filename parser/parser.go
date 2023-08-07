package parser

import (
	"structlang/ast"
	"structlang/lexer"
	"structlang/token"
)

type (
	prefixFn func() ast.Expr
	infixFn  func(left ast.Expr) ast.Expr
)

type Parser struct {
	l              *lexer.Lexer
	curr, peek     token.Token
	prefixParseFns map[token.TokenType]prefixFn
	infixParseFns  map[token.TokenType]infixFn
	errors         []error
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.getNextToken()
	p.getNextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixFn)
	p.infixParseFns = make(map[token.TokenType]infixFn)

	p.registerPrefix(token.IDENT, p.ident)
	p.registerPrefix(token.INT, p.intLit)
	p.registerPrefix(token.FLOAT, p.floatLit)
	p.registerPrefix(token.STRING, p.stringLit)
	p.registerPrefix(token.TRUE, p.boolLit)
	p.registerPrefix(token.FALSE, p.boolLit)

	return p
}

func (p Parser) Errors() []error {
	return p.errors
}

func (p *Parser) getNextToken() {
	p.curr = p.peek
	p.peek = p.l.GetNextToken()
}

func (p *Parser) peekedTokenIs(tt token.TokenType) bool { return p.peek.Type == tt }
func (p *Parser) expect(tt token.TokenType) bool {
	if p.peekedTokenIs(tt) {
		p.getNextToken()
		return true
	} else {
		p.peekError(tt)
		return false
	}
}

func (p *Parser) registerPrefix(tt token.TokenType, fn prefixFn) {
	p.prefixParseFns[tt] = fn
}
func (p *Parser) registerInfix(tt token.TokenType, fn infixFn) {
	p.infixParseFns[tt] = fn
}

func (p *Parser) Program() *ast.Program {
	program := &ast.Program{}
	program.Stmts = []ast.Stmt{}

	for p.curr.Type != token.EOF {
		stmt := p.stmt()
		if stmt != nil {
			program.Stmts = append(program.Stmts, stmt)
		}
		p.getNextToken()
	}

	return program
}
