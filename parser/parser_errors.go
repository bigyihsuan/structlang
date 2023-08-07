package parser

import (
	"fmt"
	"structlang/token"
)

func (p *Parser) peekError(tt token.TokenType) {
	err := fmt.Errorf("expected next token to be `%s`, got `%s` instead", tt, p.peek.Lexeme)
	p.errors = append(p.errors, err)
}
func (p *Parser) noPrefixParseFnError(tt token.TokenType) {
	err := fmt.Errorf("no prefix parse function for %s found", tt)
	p.errors = append(p.errors, err)
}
func (p *Parser) noInfixParseFnError(tt token.TokenType) {
	err := fmt.Errorf("no infix parse function for %s found", tt)
	p.errors = append(p.errors, err)
}
