package main

import (
	"avm/ReadFile"
	"avm/shell"
	"fmt"
	"io"
	"os"
	"strings"
)

const exitFail = 1

func run(args []string, r io.Reader ,w io.Writer) error {
	// check usage are respected
	if len(args) > 2 {
		err := fmt.Errorf("too few arguments, got %d expected %d\nusage: avm [filename.avm] ", len(args)-1, 1)
		return err
	}

	// no Args start CLI mod
	if len(args) == 1 {
		fmt.Println("Welcome Abstract VM CLI")
		sh := shell.NewShell()
		sh.Run(r, w)
	} else { // parse .avm file
		filename := os.Args[1]
		// make sur we have a .avm file as input file
		if !strings.HasSuffix(filename, ".avm") {
			ext := strings.Split(filename, ".")
			return fmt.Errorf("bad file format, got \".%s\" format but expected .avm format", ext[1])
		}

		return ReadFile.ReadFile(filename)
	}

	return nil
}

func main() {
	if err := run(os.Args, os.Stdin, os.Stdout); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFail)
	}
}
