package main

import (
	"fmt"
	"os"

	"github.com/kr/pretty"

	"github.com/bigyihsuan/structlang/eval"
	"github.com/bigyihsuan/structlang/lexer"
	"github.com/bigyihsuan/structlang/parser"
	"github.com/bigyihsuan/structlang/token"
)

func main() {
	// b, _ := os.ReadFile("example/tree.struct")
	b, _ := os.ReadFile("example/expr.struct")
	src := string(b)
	// src := `type Tree[T] = struct[T]{v T; l,r Either[Tree[T],nil] };`
	lex, _ := lexer.NewLexer(src)
	fmt.Println(src)

	var tokens []token.Token
	tok := lex.Lex()
	for tok.Type() != token.EOF {
		tokens = append(tokens, tok)
		tok = lex.Lex()
	}

	tokens = lexer.ClearComments(tokens)

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
	astparser := parser.NewAstParser(tree)
	asttree := astparser.Parse()
	// pretty.Println(asttree)

	evaluator := eval.NewEvaluator(asttree)
	evaluator.Stmt(&evaluator.BaseEnv)
	pretty.Println(evaluator.BaseEnv)
}
