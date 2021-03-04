package main

import (
	"avm/reader"
	"avm/cmd/avm/shell"
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"os"
	"strings"
)

func run(args []string, r io.Reader, w io.Writer) error {
	// check usage are respected
	if len(args) > 2 {
		return fmt.Errorf("too few arguments, got %d expected %d\nusage: avm [filename.avm] ", len(args)-1, 1)
	}
	
	// no Args start CLI mod
	if len(args) == 1 {
		app := cli.App{
			Name:                 "Abstact VM",
			Usage:                "Enter an instruction",
			EnableBashCompletion: true,
			Action: func(ctx *cli.Context) error {
				return shell.Run(os.Stdin, os.Stdout)
			},
		}
		
		if err := app.Run(os.Args); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}
	// parse .avm file
	filename := os.Args[1]
	// make sur we have a .avm file as input file
	if !strings.HasSuffix(filename, ".avm") {
		ext := strings.Split(filename, ".")
		return fmt.Errorf("bad file format, got \".%s\" format but expected .avm format", ext[1])
	}
	
	return reader.ReadFile(filename)
	
}

func main() {
	if err := run(os.Args, os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
	
}
