package main

import (
	"fmt"
	"os"

	"github.com/bigyihsuan/structlang/lexer"
	"github.com/bigyihsuan/structlang/token"
)

func main() {
	b, _ := os.ReadFile("tree.struct")
	src := string(b)
	lex, _ := lexer.NewLexer(src)

	var tokens []token.Token
	tok := lex.Lex()
	tokens = append(tokens, tok)
	for tok.Type() != token.EOF {
		tok = lex.Lex()
		tokens = append(tokens, tok)
	}
	fmt.Println(tokens)
}
