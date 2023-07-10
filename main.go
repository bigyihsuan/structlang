package main

import (
	"fmt"
	"os"

	"github.com/bigyihsuan/structlang/eval"
	"github.com/bigyihsuan/structlang/lexer"
	"github.com/bigyihsuan/structlang/parser"
	"github.com/bigyihsuan/structlang/token"
	"github.com/kr/pretty"
)

const srcTemplate = `"""
%s"""
`

func main() {
	// b, _ := os.ReadFile("example/tree.struct")
	b, _ := os.ReadFile("example/expr.struct")
	src := string(b) + "\n"
	// src := `type Tree[T] = struct[T]{v T; l,r Either[Tree[T],nil] };`
	lex, _ := lexer.NewLexer(src)
	fmt.Printf(srcTemplate, src)
	fmt.Println()

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
	// pretty.Println(tree)
	// fmt.Println()

	astparser := parser.NewAstParser(tree)
	asttree := astparser.Parse()
	pretty.Println(asttree)
	fmt.Println()

	evaluator := eval.NewEvaluator(asttree)
	err := evaluator.Stmt(&evaluator.BaseEnv)
	if err != nil {
		fmt.Println(err)
	}
	// pretty.Println(evaluator.BaseEnv)

	for id, ty := range evaluator.BaseEnv.Types {
		fmt.Printf("%s = %s\n", id.String(), ty.String())
	}
	fmt.Println()
	for id, val := range evaluator.BaseEnv.Variables {
		fmt.Printf("%s = %v\n", id.String(), val)
	}
}
