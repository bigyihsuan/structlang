package token

import (
	"fmt"
	gotoken "go/token"
)

type Token struct {
	Type     TokenType
	Lexeme   string
	Position gotoken.Position
}

func New(tokenType TokenType, lexeme, filename string, offset, line, column int) Token {
	t := Token{
		Type:     tokenType,
		Lexeme:   lexeme,
		Position: gotoken.Position{Filename: filename, Offset: offset, Line: line, Column: column},
	}
	return t
}

func (t Token) String() string {
	return fmt.Sprintf("{%s, `%s`, %v}", t.Type, t.Lexeme, t.Position)
}
