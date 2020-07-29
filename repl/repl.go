package repl

import (
	"avm/evaluator"
	"avm/lexer"
	"avm/parser"
	"bufio"
	"fmt"
	"io"
)

const PROMPT = "avm>"

func printParserErrors(out io.Writer, errors []string) {
	for _, err := range errors {
		_, _ = io.WriteString(out, err)
	}
}

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	st := evaluator.New()
	for {
		_, _ = fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		pg := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
		}

		_, err := st.Eval(pg)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			st.Print()
		}

	}
}
