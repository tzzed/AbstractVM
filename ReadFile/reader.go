package ReadFile

import (
	"avm/evaluator"
	"avm/lexer"
	"avm/parser"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	fmt.Printf("File %s is opened...\n", filename)
	scanner := bufio.NewScanner(f)
	st := evaluator.NewStack()
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		l := lexer.New(line)
		p := parser.New(l)
		pg, pErr := p.ParseProgram()
		if pErr != nil {
			fmt.Println(pErr.Error())
			continue
		}

		if pg == nil {
			continue
		}

		_, err = st.Eval(pg)
		if err != nil {
			fmt.Println(err.Error() + line)
			return err
		} else {
			st.Dump()
		}
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	return nil
}
