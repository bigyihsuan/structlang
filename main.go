package main

import (
	"fmt"
	"os"

	"structlang/lexer"
	"structlang/repl"
	"structlang/token"

	"github.com/jessevdk/go-flags"
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

	if opts.File == "" && opts.Code == "" {
		repl.Start(os.Stdin, os.Stdout)
	} else if opts.File != "" && opts.Code != "" {
		fmt.Fprintln(os.Stderr, "-f/--file and -c/--code flags are mutually exclusive")
		os.Exit(1)
	}

	var src string
	var filename string

	if opts.File != "" {
		bytes, err := os.ReadFile(string(opts.File))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		src = string(bytes)
		filename = string(opts.File)
	} else if opts.Code != "" {
		src = string(opts.Code)
		filename = ""
	}

	lex := lexer.New(src, filename)
	if opts.Debug {
		fmt.Printf(srcTemplate, src)
		fmt.Println()
	}

	for tok := lex.NextToken(); tok.Type != token.EOF; tok = lex.NextToken() {
		fmt.Println(tok)
	}
}
