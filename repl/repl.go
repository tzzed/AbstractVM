package repl

import (
	"avm/lexer"
	"avm/token"
	"bufio"
	"fmt"
	"io"
)

const PROMPT = "avm>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		_, _ = fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			_, _ = fmt.Fprintf(out, "{Type: %s, Literal:%s}\n", tok.Type, tok.Literal)
		}
	}
}
