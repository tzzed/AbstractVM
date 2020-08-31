package ReadFile

import (
	"avm/evaluator"
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

	scanner := bufio.NewScanner(f)
	st := evaluator.NewStack()
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		p := parser.NewParser(line)
		pg, pErr := p.ParseInstruction()
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
		}
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	_, _ = st.Dump()

	return nil
}
