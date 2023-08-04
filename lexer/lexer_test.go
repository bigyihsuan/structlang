package lexer

import (
	"testing"

	"structlang/token"
)

func TestLexer(t *testing.T) {
	tests := `let x = 10;
	12.34
	true
	false
	"hello\nworld"
	set x = 1;`

	expected := []struct {
		tokType token.TokenType
		lexeme  string
	}{
		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.EQ, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.FLOAT, "12.34"},
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.STRING, "hello\nworld"},
		{token.SET, "set"},
		{token.IDENT, "x"},
		{token.EQ, "="},
		{token.INT, "1"},
	}

	l := New(tests, "TEST")
	for i, tt := range expected {
		tok := l.NextToken()
		if tok.Type != tt.tokType {
			t.Fatalf("tests[%d] - incorrect token type: expect=%q, got=%q", i, tt.tokType, tok.Type)
		}
		if tok.Lexeme != tt.lexeme {
			t.Fatalf("tests[%d] - incorrect token lexeme: expect=%q, got=%q", i, tt.lexeme, tok.Lexeme)
		}
	}
}
