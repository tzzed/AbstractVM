package main

import (
	"avm/evaluator"
	"avm/lexer"
	"avm/parser"
	"avm/shell"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	args := os.Args
	if len(args) > 2 {
		fmt.Printf("too few arguments, got %d expected %d", len(args)-1, 1)
		fmt.Printf("usage: avm [filename]")
		return
	}

	if len(args) == 1 {
		fmt.Println("This Abstract VM CLI")
		shell.Start(os.Stdin, os.Stdout)
	} else {
		filename := os.Args[1]
		f, err := os.Open(filename)
		if err != nil {
			return
		}

		defer f.Close()

		fmt.Printf("File %s Opened...\n", filename)
		scanner := bufio.NewScanner(f)
		st := evaluator.NewStack()
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimSpace(line)
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
				return
			} else {
				st.Dump()
			}
		}

		if err = scanner.Err(); err != nil {
			log.Fatal(err)
			return
		}
	}
}
