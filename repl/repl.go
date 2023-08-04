package repl

import (
	"bufio"
	"fmt"
	"io"

	"structlang/lexer"
	"structlang/token"
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
		for t := l.NextToken(); t.Type != token.EOF; t = l.NextToken() {
			io.WriteString(out, t.String()+"\n")
		}
	}
}
