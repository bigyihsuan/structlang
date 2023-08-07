package parser

import (
	"structlang/ast"
	"structlang/token"
)

func (p *Parser) stmt() ast.Stmt {
	switch p.curr.Type {
	// case token.LET:
	// 	return p.letStmt()
	default:
		return p.exprStmt()
	}
}

func (p *Parser) exprStmt() *ast.ExprStmt {
	stmt := &ast.ExprStmt{FirstToken: p.curr}
	stmt.Expr = p.expr(LOWEST)
	if p.peekedTokenIs(token.SEMICOLON) {
		p.getNextToken()
	}
	return stmt
}

// func (p *Parser) letStmt() ast.Stmt {
// 	stmt := &ast.LetStmt{FirstToken: p.curr}
// 	if !p.expect(token.IDENT) {
// 		return nil
// 	}
// 	stmt.Name = p.lvalue()
// 	if !p.expect(token.EQ) {
// 		return nil
// 	}
// 	p.getNextToken()
// 	stmt.Value = p.expr(LOWEST)
// 	if p.peekedTokenIs(token.SEMICOLON) {
// 		p.getNextToken()
// 	}
// 	_, stmt.LastToken = stmt.Value.BoundingTokens()
// 	return stmt
// }
