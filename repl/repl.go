package repl

import (
	"bufio"
	"fmt"
	"io"

	"structlang/lexer"
	"structlang/parser"
)

const PROMPT = ">>> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New([]rune(line), "REPL")
		// for t := l.GetNextToken(); t.Type != token.EOF; t = l.GetNextToken() {
		// 	io.WriteString(out, t.String()+"\n")
		// }
		p := parser.New(l)
		program := p.Program()
		if len(p.Errors()) != 0 {
			for _, err := range p.Errors() {
				fmt.Fprint(out, err)
			}
			continue
		}
		fmt.Fprint(out, program.String()+"\n")
	}
}
