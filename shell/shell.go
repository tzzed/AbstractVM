package shell

import (
	"avm/evaluator"
	"avm/lexer"
	"avm/parser"
	"bufio"
	"fmt"
	"io"
)

const PROMPT = "avm>"

type Shell struct {
	Prompt string
}

func (sh *Shell) signalsManager() {

}

// Run shell
func (sh *Shell) Run(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	st := evaluator.NewStack()
	for {
		_, _ = fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		pg, err := p.ParseProgram()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if pg == nil {
			continue
		}

		_, err = st.Eval(pg)
		if err != nil {
			fmt.Println(err.Error() + line)

		} else {
			st.Dump()
		}

	}
}


func NewShell() *Shell {
	return &Shell{Prompt: PROMPT}
}