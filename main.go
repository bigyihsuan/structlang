package main

import (
	"fmt"
	"os"

	"github.com/bigyihsuan/structlang/lexer"
	"github.com/bigyihsuan/structlang/parser"
	"github.com/bigyihsuan/structlang/token"

	"github.com/kr/pretty"
)

func main() {
	// b, _ := os.ReadFile("tree.struct")
	// src := string(b)
	src := `type Tree[T] = struct[T]{v T; l,r Either[Tree[T],nil] };`
	lex, _ := lexer.NewLexer(src)
	fmt.Println(src)

	var tokens []token.Token
	tok := lex.Lex()
	for tok.Type() != token.EOF {
		tokens = append(tokens, tok)
		tok = lex.Lex()
	}

	for _, tok := range tokens {
		fmt.Println(tok.String())
	}
	fmt.Println()

	p := parser.NewParser(tokens)
	tree, errs := p.Parse()
	if errs != nil {
		fmt.Fprintln(os.Stderr, errs)
		return
	}
	// pretty.Println(tree)
	astparser := parser.NewAstParser(tree)
	ast := astparser.Parse()
	pretty.Println(ast)
}
