package parser

import (
	"strconv"
	"structlang/ast"
)

func (p *Parser) ident() ast.Expr {
	return &ast.Ident{
		Token: p.curr,
		Value: p.curr.Lexeme,
	}
}

func (p *Parser) intLit() ast.Expr {
	v, err := strconv.Atoi(p.curr.Lexeme)
	if err != nil {
		p.errors = append(p.errors, err)
		return nil
	}
	return &ast.IntLit{
		Token: p.curr,
		Value: int64(v),
	}
}

func (p *Parser) floatLit() ast.Expr {
	v, err := strconv.ParseFloat(p.curr.Lexeme, 64)
	if err != nil {
		p.errors = append(p.errors, err)
		return nil
	}
	return &ast.FloatLit{
		Token: p.curr,
		Value: v,
	}
}

func (p *Parser) stringLit() ast.Expr {
	return &ast.StrLit{
		Token: p.curr,
		Value: p.curr.Lexeme,
	}
}

func (p *Parser) boolLit() ast.Expr {
	return &ast.BoolLit{
		Token: p.curr,
		Value: p.curr.Lexeme == "true",
	}
}
