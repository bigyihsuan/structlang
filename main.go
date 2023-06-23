package main

import (
	"fmt"
	"os"

	"github.com/bigyihsuan/structlang/eval"
	"github.com/bigyihsuan/structlang/lexer"
	"github.com/bigyihsuan/structlang/parser"
	"github.com/bigyihsuan/structlang/token"
	"github.com/bigyihsuan/structlang/trees/ast"
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
	asttree := astparser.Parse()
	for _, node := range asttree {
		switch node := node.(type) {
		case ast.TypeDef:
			fmt.Println(node.FirstToken, node.LastToken)
		}
	}

	evaluator := eval.NewEvaluator(asttree)
	evaluator.Stmt(&evaluator.BaseEnv)
	fmt.Println(evaluator.BaseEnv)
}
