package parser

import (
	"structlang/ast"
	"structlang/lexer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExprStmt(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`x;`, "(x)"},
		{`10;`, "(10)"},
		{`1.23;`, "(1.23)"},
		{`"hello world";`, "(hello world)"},
		{`true;`, "(true)"},
		{`false;`, "(false)"},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input, "")
		p := New(l)
		program := p.Program()
		if !testProgram(t, program, 1) {
			return
		}
		stmt, ok := program.Stmts[0].(*ast.ExprStmt)
		if !ok {
			wrongTypeFatal(t, "ExprStmt", program.Stmts[0])
		}
		assert.Equal(t, tt.expected, stmt.Expr.String())
	}
}

func TestLetStmt(t *testing.T) {
	input := `let x = 10;`
	l := lexer.New(input, "")
	p := New(l)
	program := p.Program()
	if !testProgram(t, program, 1) {
		return
	}
	stmt, ok := program.Stmts[0].(*ast.LetStmt)
	if !ok {
		wrongTypeFatal(t, "LetStmt", program.Stmts[0])
	}
	name := "x"
	value := "10"
	assert.Equal(t, name, stmt.Ident.String())
	assert.Equal(t, value, stmt.Value.String())
}
