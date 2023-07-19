package main

import (
	"fmt"
	"os"

	"github.com/bigyihsuan/structlang/eval"
	"github.com/bigyihsuan/structlang/lexer"
	"github.com/bigyihsuan/structlang/parser"
	"github.com/bigyihsuan/structlang/token"
	"github.com/jessevdk/go-flags"
	"github.com/kr/pretty"
)

const srcTemplate = `"""
%s"""
`

func main() {
	var opts struct {
		File  flags.Filename `short:"f" long:"file" value-name:"FILE" description:"Input code file."`
		Code  flags.Filename `short:"c" long:"code" value-name:"CODE" description:"Argument-provided code."`
		Debug bool           `short:"d" long:"debug" description:"Output debugging information."`
	}
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	if opts.File != "" && opts.Code != "" {
		fmt.Fprintln(os.Stderr, "-f/--file and -c/--code flags are mutually exclusive")
		os.Exit(1)
	}

	var src string

	if opts.File != "" {
		bytes, err := os.ReadFile(string(opts.File))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		src = string(bytes)
	} else if opts.Code != "" {
		src = string(opts.Code)
	}
	src += "\n"

	lex, _ := lexer.NewLexer(src)
	if opts.Debug {
		fmt.Printf(srcTemplate, src)
		fmt.Println()
	}

	var tokens []token.Token
	tok := lex.Lex()
	for tok.Type() != token.EOF {
		tokens = append(tokens, tok)
		tok = lex.Lex()
	}

	tokens = lexer.ClearComments(tokens)

	if opts.Debug {
		for _, tok := range tokens {
			fmt.Println(tok.String())
		}
		fmt.Println()
	}

	p := parser.NewParser(tokens)
	tree, errs := p.Parse()
	if errs != nil {
		fmt.Fprintln(os.Stderr, errs)
		return
	}
	if opts.Debug {
		pretty.Println(tree)
		fmt.Println()
		for _, stmt := range tree {
			fmt.Println(stmt)
		}
	}

	astparser := parser.NewAstParser(tree)
	asttree := astparser.Parse()
	if opts.Debug {
		pretty.Println(asttree)
		fmt.Println()
	}

	evaluator := eval.NewEvaluator(asttree)
	err = evaluator.Evaluate(&evaluator.BaseEnv)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	// pretty.Println(evaluator.BaseEnv)
	if opts.Debug {
		fmt.Println("types:\n=======")
		for id, ty := range evaluator.BaseEnv.Types {
			fmt.Printf("%s = %s\n", id, ty.String())
		}
		fmt.Println()
		fmt.Println("vars:\n=======")
		for id, val := range evaluator.BaseEnv.Variables {
			fmt.Printf("%s %s = %v\n", id, val.TypeName(), val)
		}
	}
}
