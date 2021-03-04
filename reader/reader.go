package reader

import (
	"avm/evaluator"
	"avm/parser"
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ReadFile read instructions from a file.
func ReadFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	
	st := evaluator.NewStack()
	
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		p := parser.NewParser(line)
		pg, pErr := p.ParseInstruction()
		if pErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, pErr.Error())
			continue
		}

		if pg == nil {
			continue
		}
		
		if _, err = st.Eval(pg); err != nil {
			return err
		}
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	_, _ = st.Dump()

	return nil
}
