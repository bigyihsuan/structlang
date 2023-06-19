package token

import (
	"fmt"
	gotoken "go/token"
)

type Token struct {
	type_    TokenType
	lexeme   string
	position gotoken.Position
}

func NewToken(tokenType TokenType, lexeme string, offset, line, column int) Token {
	t := Token{
		type_:    tokenType,
		lexeme:   lexeme,
		position: gotoken.Position{Filename: "", Offset: offset, Line: line, Column: column},
	}
	return t
}

func (t Token) Type() TokenType            { return t.type_ }
func (t Token) Lexeme() string             { return t.lexeme }
func (t Token) Position() gotoken.Position { return t.position }

func (t Token) String() string {
	return fmt.Sprintf("{%s, `%s`, %v}", t.type_, t.lexeme, t.position)
}
